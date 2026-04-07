# TypeScript Types/Interfaces for Frontend

This file contains TypeScript interfaces and types for the API responses. Copy and use these in your frontend project.

## Base Types

```typescript
// Standard API Response (Single Data)
export interface ApiResponse<T> {
  success: boolean;
  statusCode: number;
  message: string;
  data: T;
}

// Standard API Response (With Pagination)
export interface ApiResponseWithPagination<T> {
  success: boolean;
  statusCode: number;
  message: string;
  data: T[];
  pagination: PaginationMeta;
}

// Pagination Metadata
export interface PaginationMeta {
  currentPage: number;
  perPage: number;
  total: number;
  lastPage: number;
  firstPageUrl: string;
  lastPageUrl: string;
  nextPageUrl?: string;
  prevPageUrl?: string;
}

// Error Response
export interface ErrorResponse {
  success: false;
  statusCode: number;
  message: string;
  error: string | null | ValidationError[];
}

// Validation Error
export interface ValidationError {
  field: string;
  message: string;
}

// HTTP Methods
export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
```

---

## Auth Types

```typescript
// User Roles
export type UserRole = 'USER' | 'ADMIN' | 'SUPER_USER';

// User Entity
export interface User {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  profileImage: string | null;
  createdAt: string; // ISO 8601 datetime
}

// Register Request
export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
  role: UserRole;
}

// Login Request
export interface LoginRequest {
  email: string;
  password: string;
}

// Refresh Token Request
export interface RefreshTokenRequest {
  refreshToken: string;
}

// Logout Request
export interface LogoutRequest {
  refreshToken: string;
}

// Auth Response (Register/Login)
export interface AuthResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

// Token Refresh Response
export interface TokenRefreshResponse {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

// API Response Types for Auth
export type RegisterResponse = ApiResponse<AuthResponse>;
export type LoginResponse = ApiResponse<AuthResponse>;
export type RefreshTokenResponse = ApiResponse<TokenRefreshResponse>;
export type LogoutResponse = ApiResponse<null>;
```

---

## User Types

```typescript
// Get Own Profile Response
export type GetOwnProfileResponse = ApiResponse<User>;

// Get All Users Response
export type GetUsersResponse = ApiResponseWithPagination<User>;

// Get User by ID Response
export type GetUserByIdResponse = ApiResponse<User>;

// Create User Request
export interface CreateUserRequest {
  name: string;
  email: string;
  password: string;
  role: UserRole;
}

// Create User Response
export type CreateUserResponse = ApiResponse<User>;

// Update User Request
export interface UpdateUserRequest {
  name?: string;
  email?: string;
  role?: UserRole;
}

// Update User Response
export type UpdateUserResponse = ApiResponse<User>;

// Delete User Response
export type DeleteUserResponse = ApiResponse<null>;

// Upload Profile Image Response
export interface UploadProfileImageResponse {
  imageUrl: string;
}

export type ProfileImageResponse = ApiResponse<UploadProfileImageResponse>;
```

---

## Employee Types

```typescript
// Employee Entity
export interface Employee {
  id: string;
  userId: string;
  name: string;
  email: string;
  phone: string;
  address: string;
  salaryBase: number;
  joinDate: string; // ISO 8601 date
  scheduleId: string | null;
  scheduleName: string | null;
  bankName: string;
  bankAccountNumber: string;
  bankAccountHolder: string;
  createdAt: string; // ISO 8601 datetime
}

// Get All Employees Response
export type GetAllEmployeesResponse = ApiResponseWithPagination<Employee>;

// Get Employee by ID Response
export type GetEmployeeByIdResponse = ApiResponse<Employee>;

// Create Employee Request
export interface CreateEmployeeRequest {
  userId: string;
  phone: string;
  address: string;
  salaryBase: number;
  joinDate: string; // YYYY-MM-DD
  scheduleId?: string;
  bankName: string;
  bankAccountNumber: string;
  bankAccountHolder: string;
}

// Create Employee Response
export type CreateEmployeeResponse = ApiResponse<Employee>;

// Update Employee Request
export interface UpdateEmployeeRequest {
  phone?: string;
  address?: string;
  salaryBase?: number;
  scheduleId?: string;
  bankName?: string;
  bankAccountNumber?: string;
  bankAccountHolder?: string;
}

// Update Employee Response
export type UpdateEmployeeResponse = ApiResponse<Employee>;

// Delete Employee Response
export type DeleteEmployeeResponse = ApiResponse<null>;
```

---

## Schedule Types

```typescript
// Schedule Entity
export interface Schedule {
  id: string;
  name: string;
  timeIn: string; // HH:MM format
  timeOut: string; // HH:MM format
  allowedLateMinutes: number;
  officeLat: number;
  officeLong: number;
  allowedRadiusMeters: number;
  createdAt: string; // ISO 8601 datetime
}

// Get All Schedules Response
export type GetAllSchedulesResponse = ApiResponseWithPagination<Schedule>;

// Get Schedule by ID Response
export type GetScheduleByIdResponse = ApiResponse<Schedule>;

// Create Schedule Request
export interface CreateScheduleRequest {
  name: string;
  timeIn: string; // HH:MM format
  timeOut: string; // HH:MM format
  allowedLateMinutes: number;
  officeLat: number;
  officeLong: number;
  allowedRadiusMeters: number;
}

// Create Schedule Response
export type CreateScheduleResponse = ApiResponse<Schedule>;

// Update Schedule Request
export interface UpdateScheduleRequest {
  name?: string;
  timeIn?: string;
  timeOut?: string;
  allowedLateMinutes?: number;
  officeLat?: number;
  officeLong?: number;
  allowedRadiusMeters?: number;
}

// Update Schedule Response
export type UpdateScheduleResponse = ApiResponse<Schedule>;

// Delete Schedule Response
export type DeleteScheduleResponse = ApiResponse<null>;
```

---

## Attendance Types

```typescript
// Attendance Status
export type AttendanceStatus = 'PRESENT' | 'LATE';

// Attendance Entity
export interface Attendance {
  id: string;
  date: string; // ISO 8601 date
  clockInTime: string | null; // ISO 8601 datetime
  clockOutTime: string | null; // ISO 8601 datetime
  status: AttendanceStatus;
  notes: string;
  scheduleName: string | null;
}

// Clock In Request
export interface ClockInRequest {
  lat: number;
  long: number;
}

// Clock In Response
export interface ClockInResponseData {
  attendanceId: string;
  employeeId: string;
  clockInTime: string;
  status: AttendanceStatus;
  distance: number;
  scheduleName: string;
}

export type ClockInResponse = ApiResponse<ClockInResponseData>;

// Clock Out Request
export interface ClockOutRequest {
  lat: number;
  long: number;
}

// Clock Out Response
export interface ClockOutResponseData {
  attendanceId: string;
  clockOutTime: string;
  distance: number;
}

export type ClockOutResponse = ApiResponse<ClockOutResponseData>;

// Get Own History Response
export type GetHistoryResponse = ApiResponseWithPagination<Attendance>;

// Get All Attendances Response (Admin)
export type GetAllAttendancesResponse = ApiResponseWithPagination<Attendance>;

// Get All Attendances Filter
export interface GetAllAttendancesFilters {
  employee_id?: string;
  schedule_id?: string;
  status?: AttendanceStatus;
  date_from?: string; // YYYY-MM-DD
  date_to?: string; // YYYY-MM-DD
}
```

---

## Payroll Types

```typescript
// Payroll Status
export type PayrollStatus = 'DRAFT' | 'APPROVED' | 'PAID';

// Payroll Item Type
export type PayrollItemType = 'EARNING' | 'DEDUCTION';

// Payroll Item
export interface PayrollItem {
  id: string;
  name: string;
  amount: number;
  type: PayrollItemType;
}

// Payroll List Item (for Get All)
export interface PayrollListItem {
  id: string;
  employeeName: string;
  period: string;
  netSalary: number;
  status: PayrollStatus;
  generatedAt: string; // ISO 8601 datetime
}

// Payroll Detail (for Get by ID)
export interface PayrollDetail {
  id: string;
  employeeId: string;
  employeeName: string;
  bankName: string;
  bankAccountNumber: string;
  bankAccountHolder: string;
  periodStart: string; // YYYY-MM-DD
  periodEnd: string; // YYYY-MM-DD
  baseSalary: number;
  totalAllowance: number;
  totalDeduction: number;
  netSalary: number;
  status: PayrollStatus;
  items: PayrollItem[];
  generatedAt: string; // ISO 8601 datetime
  createdAt: string; // ISO 8601 datetime
  updatedAt: string; // ISO 8601 datetime
}

// Generate Bulk Payroll Request
export interface GenerateBulkPayrollRequest {
  periodMonth: number; // 1-12
  periodYear: number; // 2020-2100
}

// Generate Bulk Payroll Response Data
export interface GenerateBulkPayrollResponseData {
  totalGenerated: number;
  periodStart: string;
  periodEnd: string;
  message: string;
}

export type GenerateBulkPayrollResponse = ApiResponse<GenerateBulkPayrollResponseData>;

// Get All Payrolls Response
export type GetAllPayrollsResponse = ApiResponseWithPagination<PayrollListItem>;

// Get Payroll by ID Response
export type GetPayrollByIdResponse = ApiResponse<PayrollDetail>;

// Update Payroll Status Request
export interface UpdatePayrollStatusRequest {
  status: PayrollStatus;
}

// Update Payroll Status Response
export type UpdatePayrollStatusResponse = ApiResponse<null>;

// Export CSV Filter
export interface ExportPayrollCSVFilters {
  month: number; // 1-12
  year: number; // 2020-2100
}
```

---

## API Client Helper (Example)

```typescript
// API Client Configuration
export interface ApiClientConfig {
  baseURL: string;
  getToken: () => string | null;
  onTokenExpired?: () => void;
}

// API Client Class (Example)
export class ApiClient {
  private config: ApiClientConfig;

  constructor(config: ApiClientConfig) {
    this.config = config;
  }

  private async request<T>(
    method: HttpMethod,
    endpoint: string,
    data?: any,
    useAuth: boolean = true
  ): Promise<T> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    if (useAuth) {
      const token = this.config.getToken();
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }

    const response = await fetch(`${this.config.baseURL}${endpoint}`, {
      method,
      headers,
      body: data ? JSON.stringify(data) : undefined,
    });

    const result = await response.json();

    if (!result.success) {
      throw new Error(result.message);
    }

    return result;
  }

  // Auth Methods
  async register(data: RegisterRequest): Promise<RegisterResponse> {
    return this.request<RegisterResponse>('POST', '/auth/register', data, false);
  }

  async login(data: LoginRequest): Promise<LoginResponse> {
    return this.request<LoginResponse>('POST', '/auth/login', data, false);
  }

  async refreshToken(data: RefreshTokenRequest): Promise<RefreshTokenResponse> {
    return this.request<RefreshTokenResponse>('POST', '/auth/refresh', data, false);
  }

  async logout(data: LogoutRequest): Promise<LogoutResponse> {
    return this.request<LogoutResponse>('POST', '/auth/logout', data);
  }

  async logoutAll(): Promise<LogoutResponse> {
    return this.request<LogoutResponse>('POST', '/auth/logout-all');
  }

  // User Methods
  async getOwnProfile(): Promise<GetOwnProfileResponse> {
    return this.request<GetOwnProfileResponse>('GET', '/users/me');
  }

  async getUsers(page: number = 1, perPage: number = 10): Promise<GetUsersResponse> {
    return this.request<GetUsersResponse>(
      'GET',
      `/users?page=${page}&per_page=${perPage}`
    );
  }

  async getUserById(id: string): Promise<GetUserByIdResponse> {
    return this.request<GetUserByIdResponse>('GET', `/users/${id}`);
  }

  async createUser(data: CreateUserRequest): Promise<CreateUserResponse> {
    return this.request<CreateUserResponse>('POST', '/users', data);
  }

  async updateUser(id: string, data: UpdateUserRequest): Promise<UpdateUserResponse> {
    return this.request<UpdateUserResponse>('PATCH', `/users/${id}`, data);
  }

  async deleteUser(id: string): Promise<DeleteUserResponse> {
    return this.request<DeleteUserResponse>('DELETE', `/users/${id}`);
  }

  // Employee Methods
  async getEmployees(
    page: number = 1,
    perPage: number = 10,
    search?: string
  ): Promise<GetAllEmployeesResponse> {
    const searchParam = search ? `&search=${search}` : '';
    return this.request<GetAllEmployeesResponse>(
      'GET',
      `/employees?page=${page}&per_page=${perPage}${searchParam}`
    );
  }

  async getEmployeeById(id: string): Promise<GetEmployeeByIdResponse> {
    return this.request<GetEmployeeByIdResponse>('GET', `/employees/${id}`);
  }

  async createEmployee(data: CreateEmployeeRequest): Promise<CreateEmployeeResponse> {
    return this.request<CreateEmployeeResponse>('POST', '/employees', data);
  }

  async updateEmployee(id: string, data: UpdateEmployeeRequest): Promise<UpdateEmployeeResponse> {
    return this.request<UpdateEmployeeResponse>('PATCH', `/employees/${id}`, data);
  }

  async deleteEmployee(id: string): Promise<DeleteEmployeeResponse> {
    return this.request<DeleteEmployeeResponse>('DELETE', `/employees/${id}`);
  }

  // Schedule Methods
  async getSchedules(page: number = 1, perPage: number = 10): Promise<GetAllSchedulesResponse> {
    return this.request<GetAllSchedulesResponse>(
      'GET',
      `/schedules?page=${page}&per_page=${perPage}`
    );
  }

  async getScheduleById(id: string): Promise<GetScheduleByIdResponse> {
    return this.request<GetScheduleByIdResponse>('GET', `/schedules/${id}`);
  }

  async createSchedule(data: CreateScheduleRequest): Promise<CreateScheduleResponse> {
    return this.request<CreateScheduleResponse>('POST', '/schedules', data);
  }

  async updateSchedule(id: string, data: UpdateScheduleRequest): Promise<UpdateScheduleResponse> {
    return this.request<UpdateScheduleResponse>('PATCH', `/schedules/${id}`, data);
  }

  async deleteSchedule(id: string): Promise<DeleteScheduleResponse> {
    return this.request<DeleteScheduleResponse>('DELETE', `/schedules/${id}`);
  }

  // Attendance Methods
  async clockIn(data: ClockInRequest): Promise<ClockInResponse> {
    return this.request<ClockInResponse>('POST', '/attendance/clock-in', data);
  }

  async clockOut(data: ClockOutRequest): Promise<ClockOutResponse> {
    return this.request<ClockOutResponse>('POST', '/attendance/clock-out', data);
  }

  async getHistory(page: number = 1, perPage: number = 10): Promise<GetHistoryResponse> {
    return this.request<GetHistoryResponse>(
      'GET',
      `/attendance/history?page=${page}&per_page=${perPage}`
    );
  }

  async getAllAttendances(
    page: number = 1,
    perPage: number = 10,
    filters?: GetAllAttendancesFilters
  ): Promise<GetAllAttendancesResponse> {
    let queryParams = `?page=${page}&per_page=${perPage}`;

    if (filters?.employee_id) queryParams += `&employee_id=${filters.employee_id}`;
    if (filters?.schedule_id) queryParams += `&schedule_id=${filters.schedule_id}`;
    if (filters?.status) queryParams += `&status=${filters.status}`;
    if (filters?.date_from) queryParams += `&date_from=${filters.date_from}`;
    if (filters?.date_to) queryParams += `&date_to=${filters.date_to}`;

    return this.request<GetAllAttendancesResponse>(
      'GET',
      `/attendance/all${queryParams}`
    );
  }

  // Payroll Methods
  async generateBulkPayroll(
    data: GenerateBulkPayrollRequest
  ): Promise<GenerateBulkPayrollResponse> {
    return this.request<GenerateBulkPayrollResponse>('POST', '/payrolls/generate', data);
  }

  async getPayrolls(page: number = 1, perPage: number = 15): Promise<GetAllPayrollsResponse> {
    return this.request<GetAllPayrollsResponse>(
      'GET',
      `/payrolls?page=${page}&per_page=${perPage}`
    );
  }

  async getPayrollById(id: string): Promise<GetPayrollByIdResponse> {
    return this.request<GetPayrollByIdResponse>('GET', `/payrolls/${id}`);
  }

  async updatePayrollStatus(
    id: string,
    data: UpdatePayrollStatusRequest
  ): Promise<UpdatePayrollStatusResponse> {
    return this.request<UpdatePayrollStatusResponse>('PATCH', `/payrolls/${id}/status`, data);
  }

  async exportPayrollCSV(filters: ExportPayrollCSVFilters): Promise<Blob> {
    const token = this.config.getToken();
    const response = await fetch(
      `${this.config.baseURL}/payrolls/export/csv?month=${filters.month}&year=${filters.year}`,
      {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      }
    );

    if (!response.ok) {
      throw new Error('Failed to export CSV');
    }

    return response.blob();
  }
}

// Usage Example
const apiClient = new ApiClient({
  baseURL: 'http://localhost:8080/api/v1',
  getToken: () => localStorage.getItem('accessToken'),
  onTokenExpired: () => {
    // Handle token refresh logic
  },
});

// Login
const loginResponse = await apiClient.login({
  email: 'john@example.com',
  password: 'password123',
});

// Store tokens
localStorage.setItem('accessToken', loginResponse.data.accessToken);
localStorage.setItem('refreshToken', loginResponse.data.refreshToken);

// Get users
const usersResponse = await apiClient.getUsers(1, 10);
console.log(usersResponse.data); // Array<User>
console.log(usersResponse.pagination); // PaginationMeta

// Get all attendances with filters
const attendancesResponse = await apiClient.getAllAttendances(1, 10, {
  status: 'PRESENT',
  date_from: '2024-01-01',
  date_to: '2024-01-31',
});
console.log(attendancesResponse.data); // Array<Attendance>
```

---

## React Hooks Example (Optional)

```typescript
import { useState, useEffect } from 'react';

// Custom hook for fetching paginated data
export function usePaginatedData<T>(
  fetchFunction: (page: number, perPage: number) => Promise<ApiResponseWithPagination<T>>,
  initialPage: number = 1,
  perPage: number = 10
) {
  const [data, setData] = useState<T[]>([]);
  const [pagination, setPagination] = useState<PaginationMeta | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchData = async (page: number = initialPage) => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetchFunction(page, perPage);
      setData(response.data);
      setPagination(response.pagination);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  return {
    data,
    pagination,
    loading,
    error,
    refetch: fetchData,
  };
}

// Usage
function UsersList() {
  const { data: users, pagination, loading, error, refetch } = usePaginatedData<User>(
    (page, perPage) => apiClient.getUsers(page, perPage)
  );

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>Users</h1>
      <ul>
        {users.map((user) => (
          <li key={user.id}>{user.name}</li>
        ))}
      </ul>
      <div>
        <button onClick={() => refetch(pagination?.currentPage - 1 || 1)}>Previous</button>
        <span>Page {pagination?.currentPage} of {pagination?.lastPage}</span>
        <button onClick={() => refetch((pagination?.currentPage || 0) + 1)}>Next</button>
      </div>
    </div>
  );
}
```

---

## Vue Composables Example (Optional)

```typescript
import { ref, Ref } from 'vue';

export function usePaginatedData<T>(
  fetchFunction: (page: number, perPage: number) => Promise<ApiResponseWithPagination<T>>,
  initialPage: number = 1,
  perPage: number = 10
) {
  const data: Ref<T[]> = ref([]);
  const pagination: Ref<PaginationMeta | null> = ref(null);
  const loading = ref(false);
  const error: Ref<string | null> = ref(null);

  const fetchData = async (page: number = initialPage) => {
    loading.value = true;
    error.value = null;
    try {
      const response = await fetchFunction(page, perPage);
      data.value = response.data;
      pagination.value = response.pagination;
    } catch (err: any) {
      error.value = err.message;
    } finally {
      loading.value = false;
    }
  };

  fetchData();

  return {
    data,
    pagination,
    loading,
    error,
    refetch: fetchData,
  };
}
```

---

## Notes

1. **Date Format**: All dates are in ISO 8601 format (e.g., `2024-01-30T10:00:00Z`)
2. **Pagination**: Most list endpoints support pagination with `page` and `per_page` parameters
3. **Authentication**: Most endpoints require Bearer token in Authorization header
4. **Error Handling**: Always check `success` field in response before accessing `data`
5. **File Upload**: Use `multipart/form-data` for profile image uploads
6. **CSV Download**: Payroll export returns raw CSV file, not JSON
