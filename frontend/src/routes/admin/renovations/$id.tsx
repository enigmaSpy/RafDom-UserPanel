import { useState, useEffect } from 'react';
import { createFileRoute, useNavigate, Link } from '@tanstack/react-router';
import { apiClient } from '../../../services/api';
import { CostEstimateTab } from '../../../views/Renovation/CostEstimateTab';
import { ChatTab } from '../../../views/Renovation/ChatTab';
export const Route = createFileRoute('/admin/renovations/$id')({
  component: RenovationDetailsPage,
});

interface ClientInfo {
  id: string;
  name: string;
  email: string;
  phone: string;
  address: string;
  city: string;
}
interface FinancialSummary {
  total_labor: number;
  labor_paid: number;
  labor_balance: number;
  material_deposits: number;
  material_expenses: number;
  deposit_balance: number;
}

interface RenovationDetailUI {
  id: string;
  name: string;
  description: string;
  status: string;
  created_at: string;
  client: ClientInfo;
  finances: FinancialSummary;
}

function RenovationDetailsPage() {
  const { id } = Route.useParams();
  const navigate = useNavigate();
  
  const [data, setData] = useState<RenovationDetailUI | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(false);
  const [activeTab, setActiveTab] = useState<'overview' | 'estimate' | 'transactions' | 'progress' | 'chat'>('overview');

  useEffect(() => {
  const fetchDetails = async () => {
    try {
      const [renovationRes, summaryRes] = await Promise.all([
        apiClient.get(`/api/renovations/${id}`),
        apiClient.get(`/api/renovations/${id}/summary`)
      ]);
      
      const rawData = renovationRes.data;
      const financeData = summaryRes.data;

      if (!rawData || !financeData) throw new Error("Brak kompletnych danych");

      setData({
        id: rawData.ID,
        name: rawData.Name,
        description: rawData.Description,
        status: rawData.Status,
        created_at: rawData.CreatedAt ? rawData.CreatedAt.split('T')[0] : '-',
        client: {
          id: rawData.Client?.ID || '',
          name: `${rawData.Client?.Name || ''} ${rawData.Client?.Surname || ''}`.trim(),
          email: rawData.Client?.Email || '-',
          phone: rawData.Client?.Phone || '-',
          address: rawData.Client?.Address || '-',
          city: rawData.Client?.City || '-'
        },
        finances: {
          total_labor: financeData.total_labor_cost,
          labor_paid: financeData.labor_paid,
          labor_balance: financeData.labor_balance,
          material_deposits: financeData.material_deposit,
          material_expenses: financeData.material_expenses,
          deposit_balance: financeData.material_balance
        }
      });
    } catch (err) {
      console.error("Błąd pobierania szczegółów", err);
      setError(true);
    } finally {
      setIsLoading(false);
    }
  };
  
  fetchDetails();
}, [id]);

  if (isLoading) return <div className="p-6 text-text-main animate-pulse">Nawiązywanie połączenia z bazą danych...</div>;
  if (error || !data) return <div className="p-6 text-red-500 bg-red-900/20 border border-red-500 rounded">Błąd krytyczny: Nie udało się załadować danych remontu.</div>;

  const formatPLN = (amount: number) => `${amount.toLocaleString('pl-PL', { minimumFractionDigits: 2 })} zł`;
  return (
    <div className="p-6 font-main max-w-7xl mx-auto text-text-main">
      
     
      <div className="mb-6 flex items-center gap-4">
        <button 
          onClick={() => navigate({ to: '/' })}
          className="p-2 hover:bg-bg-surface rounded transition-colors text-text-muted hover:text-text-main"
        >
          &larr; Wróć
        </button>
        <div>
          <div className="flex items-center gap-3">
            <h1 className="text-2xl font-bold">{data.name}</h1>
            <span className="bg-blue-900/50 text-blue-400 border border-blue-800 px-2 py-0.5 rounded text-xs font-medium uppercase">
              {data.status.replace('_', ' ')}
            </span>
          </div>
          <p className="text-sm text-text-muted">PRJ-{data.id.substring(0, 4).toUpperCase()} · {data.client.name}</p>
        </div>
      </div>

   
      <div className="flex border-b border-bg-border mb-6 overflow-x-auto">
        {[
          { id: 'overview', label: 'Overview' },
          { id: 'estimate', label: 'Cost Estimate' },
          { id: 'transactions', label: 'Transactions' },
          { id: 'progress', label: 'Progress Log' },
          { id: 'chat', label: 'Chat' },
        ].map(tab => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id as any)}
            className={`px-6 py-3 text-sm font-medium transition-colors border-b-2 whitespace-nowrap ${
              activeTab === tab.id 
                ? 'border-primary text-primary' 
                : 'border-transparent text-text-muted hover:text-text-main hover:border-bg-border'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      
      {activeTab === 'overview' && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          
        
          <div className="space-y-6">
            
            
            <div className="bg-bg-surface border border-bg-border rounded-card p-6">
              <h2 className="text-xs font-bold text-text-muted uppercase tracking-wider mb-6">Project Information</h2>
              
              <div className="space-y-4">
                <div>
                  <div className="text-xs text-text-muted mb-1">Project Name</div>
                  <div className="font-medium">{data.name}</div>
                </div>
                <div>
                  <div className="text-xs text-text-muted mb-1">Description</div>
                  <div className="text-sm text-text-muted leading-relaxed">{data.description}</div>
                </div>
                <div>
                  <div className="text-xs text-text-muted mb-1">Status</div>
                  
                  <select 
                    className="bg-bg-base border border-bg-border text-sm p-2 rounded focus:outline-none focus:border-primary"
                    defaultValue={data.status}
                  >
                    <option value="estimation">Estimation</option>
                    <option value="in_progress">In Progress</option>
                    <option value="completed">Completed</option>
                  </select>
                </div>
                <div>
                  <div className="text-xs text-text-muted mb-1">Created</div>
                  <div className="text-sm">{data.created_at}</div>
                </div>
              </div>
            </div>

            <div className="bg-bg-surface border border-bg-border rounded-card p-6">
              <h2 className="text-xs font-bold text-text-muted uppercase tracking-wider mb-6">Client</h2>
              
              <div className="flex items-center gap-4 mb-4">
                <div className="w-12 h-12 bg-primary/20 text-primary font-bold flex items-center justify-center rounded-full">
                  {data.client.name.split(' ').map(n => n[0]).join('')}
                </div>
                <div>
                  <div className="font-bold">{data.client.name}</div>
                  <div className="text-xs text-text-muted">{data.client.city}</div>
                </div>
              </div>

              <div className="space-y-2 text-sm text-text-muted">
                <div className="flex gap-2"><span>✉</span> {data.client.email}</div>
                <div className="flex gap-2"><span>📞</span> {data.client.phone}</div>
                <div className="flex gap-2"><span>📍</span> {data.client.address}, {data.client.city}</div>
              </div>
            </div>

          </div>

          <div>
            <div className="bg-bg-surface border border-bg-border rounded-card p-6">
              <h2 className="text-xs font-bold text-text-muted uppercase tracking-wider mb-6">Financial Summary</h2>
              
              <div className="space-y-4 text-sm">
                
                <div className="flex justify-between items-center py-2 border-b border-bg-border">
                  <div className="text-text-muted">Total Labor<br/><span className="text-xs opacity-50">5 tasks</span></div>
                  <div className="font-bold">{formatPLN(data.finances.total_labor)}</div>
                </div>
                <div className="flex justify-between items-center py-2 border-b border-bg-border">
                  <div className="text-text-muted">Labor Paid<br/><span className="text-xs opacity-50">via payments</span></div>
                  <div className="font-medium text-emerald-400">{formatPLN(data.finances.labor_paid)}</div>
                </div>
                <div className="flex justify-between items-center py-2 border-b border-bg-border">
                  <div className="text-text-muted">Labor Balance<br/><span className="text-xs opacity-50">outstanding</span></div>
                  <div className="font-medium text-red-400">{formatPLN(data.finances.labor_balance)}</div>
                </div>

                <div className="flex justify-between items-center py-2 pt-4">
                  <div className="text-text-muted">Material Deposits</div>
                  <div className="font-medium">{formatPLN(data.finances.material_deposits)}</div>
                </div>
                <div className="flex justify-between items-center py-2 border-b border-bg-border">
                  <div className="text-text-muted">Material Expenses</div>
                  <div className="font-medium text-red-400">-{formatPLN(data.finances.material_expenses)}</div>
                </div>

                <div className="flex justify-between items-center p-3 bg-bg-base rounded font-bold mt-4">
                  <div>Deposit Balance</div>
                  <div className="text-emerald-400">{formatPLN(data.finances.deposit_balance)}</div>
                </div>

              </div>
            </div>
          </div>

        </div>
      )}

      {activeTab === 'estimate' && <CostEstimateTab renovationId={id} />}
      {activeTab === 'transactions' && <div className="p-6 bg-bg-surface rounded border border-bg-border">Sekcja Transakcji w budowie...</div>}
      {activeTab === 'progress' && <div className="p-6 bg-bg-surface rounded border border-bg-border">Sekcja Dziennika Prac w budowie...</div>}
      {activeTab === 'chat' && <ChatTab renovationId={id} />}

    </div>
  );
}