export const Tile = ({ title, value, highlight = false }: { title: string, value?: string, highlight?: boolean }) => (
    <div className="bg-bg-surface border border-bg-border p-6 rounded-[12px] shadow-sm hover:shadow-md transition-shadow">
        <h3 className="text-sm font-medium text-text-muted uppercase tracking-wider mb-2">
            {title}
        </h3>
        <div className={`text-3xl font-bold ${highlight ? 'text-primary' : 'text-text-main'}`}>
            {value || '0'}
        </div>
    </div>
);