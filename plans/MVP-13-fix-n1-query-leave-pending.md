# MVP-13: Fix N+1 Query in Leave GetPendingRequests

## Prioritas: 🟡 IMPORTANT — Performance
## Estimasi: 2 jam
## Tipe: Performance Optimization

---

## Deskripsi Masalah

Di `internal/leave/service/leave_service.go` line 255-276, `GetPendingLeaveRequests`:
```go
for i, req := range requests {
    user, _ := s.userRepo.FindByID(ctx, req.EmployeeID)     // N+1 (MongoDB)
    employee, _ := s.employeeRepo.FindByID(ctx, ...)        // N+1 (PostgreSQL)
    leaveType, _ := s.leaveTypeRepo.FindByID(ctx, ...)      // N+1 (PostgreSQL)
}
```

10 pending requests = **31 queries** (1 FindPending + 3×10 individual lookups). Ini sangat lambat dan mempengaruhi response time halaman approval admin.

## Solusi

1. Batch fetch employees by IDs (1 query)
2. Batch fetch leave types by IDs (1 query)
3. Batch fetch users by IDs dari MongoDB (1 query)
4. Map ke responses menggunakan in-memory lookup

## File yang Diubah

### 1. [MODIFY] `internal/user/repository/user_repository.go`

**Tambah method di interface:**
```go
FindByIDs(ctx context.Context, ids []string) ([]*entity.User, error)
```

### 2. [MODIFY] Implementasi user repository (MongoDB)

```go
func (r *userRepositoryImpl) FindByIDs(ctx context.Context, ids []string) ([]*entity.User, error) {
    objectIDs := make([]primitive.ObjectID, 0, len(ids))
    for _, id := range ids {
        oid, err := primitive.ObjectIDFromHex(id)
        if err == nil {
            objectIDs = append(objectIDs, oid)
        }
    }

    filter := bson.M{"_id": bson.M{"$in": objectIDs}}
    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var users []*entity.User
    if err := cursor.All(ctx, &users); err != nil {
        return nil, err
    }
    return users, nil
}
```

### 3. [MODIFY] `internal/leave/repository/leave_type_repository.go`

**Tambah method di interface:**
```go
FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*entity.LeaveType, error)
```

**Implementasi:**
```go
func (r *leaveTypeRepositoryImpl) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*entity.LeaveType, error) {
    query := `SELECT id, name, code, description, max_days, is_paid, is_active, created_at, updated_at
              FROM leave_types WHERE id = ANY($1)`
    rows, err := r.pool.Query(ctx, query, ids)
    if err != nil { return nil, err }
    defer rows.Close()

    var types []*entity.LeaveType
    for rows.Next() {
        lt := &entity.LeaveType{}
        err := rows.Scan(&lt.ID, &lt.Name, &lt.Code, &lt.Description,
            &lt.MaxDays, &lt.IsPaid, &lt.IsActive, &lt.CreatedAt, &lt.UpdatedAt)
        if err != nil { return nil, err }
        types = append(types, lt)
    }
    return types, nil
}
```

### 4. [MODIFY] `internal/leave/service/leave_service.go`

**Replace `GetPendingLeaveRequests` (line 255-276):**
```go
func (s *leaveService) GetPendingLeaveRequests(ctx context.Context) ([]dto.LeaveRequestResponse, error) {
    requests, err := s.leaveRequestRepo.FindPending(ctx)
    if err != nil { return nil, err }

    if len(requests) == 0 {
        return []dto.LeaveRequestResponse{}, nil
    }

    // Collect unique IDs
    userIDs := make(map[string]bool)
    employeeIDs := make(map[uuid.UUID]bool)
    leaveTypeIDs := make(map[uuid.UUID]bool)
    for _, req := range requests {
        userIDs[req.EmployeeID] = true
        if empID, err := uuid.Parse(req.EmployeeID); err == nil { employeeIDs[empID] = true }
        if ltID, err := uuid.Parse(req.LeaveTypeID); err == nil { leaveTypeIDs[ltID] = true }
    }

    // Batch fetch users
    userIDSlice := make([]string, 0, len(userIDs))
    for id := range userIDs { userIDSlice = append(userIDSlice, id) }
    users, _ := s.userRepo.FindByIDs(ctx, userIDSlice)
    userMap := make(map[string]*userEntity.User)
    for _, u := range users { userMap[u.ID] = u }

    // Batch fetch employees
    empIDSlice := make([]uuid.UUID, 0, len(employeeIDs))
    for id := range employeeIDs { empIDSlice = append(empIDSlice, id) }
    employees, _ := s.employeeRepo.FindByIDs(ctx, empIDSlice)
    empMap := make(map[string]*employeeEntity.Employee)
    for _, e := range employees { empMap[e.ID.String()] = e }

    // Batch fetch leave types
    ltIDSlice := make([]uuid.UUID, 0, len(leaveTypeIDs))
    for id := range leaveTypeIDs { ltIDSlice = append(ltIDSlice, id) }
    leaveTypes, _ := s.leaveTypeRepo.FindByIDs(ctx, ltIDSlice)
    ltMap := make(map[string]*entity.LeaveType)
    for _, lt := range leaveTypes { ltMap[lt.ID.String()] = lt }

    // Build responses
    responses := make([]dto.LeaveRequestResponse, 0, len(requests))
    for _, req := range requests {
        user, _ := userMap[req.EmployeeID]
        emp, _ := empMap[req.EmployeeID]
        lt, _ := ltMap[req.LeaveTypeID]

        if user == nil || emp == nil || lt == nil { continue }

        responses = append(responses, *s.toLeaveRequestResponse(
            &req, user.Name, user.Email, emp.Position,
            lt.Name, lt.Code, lt.IsPaid,
        ))
    }

    return responses, nil
}
```

> **Catatan**: Pattern yang sama bisa diterapkan ke `GetMyLeaveRequests` yang juga punya N+1 di line 222-224.

## Verifikasi

1. `go build ./...` — compile sukses
2. Buat 20+ pending leave requests
3. `GET /api/v1/leave/requests/pending` — response cepat dan benar
4. Monitor query count (via database logging): harus 4 query total (findPending + 3 batch), bukan 31+
