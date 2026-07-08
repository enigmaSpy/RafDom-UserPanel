export interface DashboardStats {
    total_projects: number;
    active_projects: number;
    completed_projects: number;
    total_clients: number;
    total_income: number;
}
export interface RenovationUI {
    id: string;
    name: string;
    client_name: string;
    status: string;
    total_budget: number;
}

export interface BackendRenovation {
    ID: string;
    Name: string;
    Status: string;
    Client: {
        Name: string;
        Surname: string;
    };
    LaborTasks: Array<{ Amount: number }> | null;
}