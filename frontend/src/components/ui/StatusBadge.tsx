export const StatusBadge = ({ status }: { status: string }) => {
    const safeStatus = status?.toLowerCase() || '';
    const statusConfig: Record<string, { label: string; styles: string }> = {
        'estimation': { label: 'Estimation', styles: 'bg-slate-700 text-slate-300' },
        'in_progress': { label: 'In Progress', styles: 'bg-blue-900/50 text-blue-400 border border-blue-800' },
        'completed': { label: 'Completed', styles: 'bg-emerald-900/50 text-emerald-400 border border-emerald-800' },
    };

    const config = statusConfig[safeStatus] || { label: safeStatus, styles: 'bg-slate-700 text-slate-300' };

    return (
        <span className={`px-3 py-1 rounded-full text-xs font-medium ${config.styles}`}>
            {config.label}
        </span>
    );
};