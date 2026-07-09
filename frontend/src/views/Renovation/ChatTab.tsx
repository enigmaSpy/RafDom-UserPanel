import { useState, useEffect, useRef } from 'react';
import { apiClient } from '../../services/api';
import { useAuth } from '../../context/AuthContext';

interface Message {
    ID: string;
    RenovationID: string;
    SenderID: string;
    ReceiverID: string;
    Content: string;
    CreatedAt: string;
}

interface ChatTabProps {
    renovationId: string;
}

export const ChatTab = ({ renovationId }: ChatTabProps) => {
    const { token, user } = useAuth();
    
    const [messages, setMessages] = useState<Message[]>([]);
    const [inputMsg, setInputMsg] = useState('');
    const [isConnected, setIsConnected] = useState(false);
    const [connectionError, setConnectionError] = useState<string | null>(null);
    
    const ws = useRef<WebSocket | null>(null);
    const messagesEndRef = useRef<HTMLDivElement>(null);

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    };

    useEffect(() => {
        scrollToBottom();
    }, [messages]);

    useEffect(() => {
        let isMounted = true;
        
        const fetchHistory = async () => {
            try {
                const res = await apiClient.get<Message[]>(`/api/renovations/${renovationId}/messages`);
                if (isMounted) {
                    setMessages(res.data || []);
                }
            } catch (err) {
                console.error("Błąd ładowania historii:", err);
            }
        };

        fetchHistory();
        
        const rawToken = localStorage.getItem('token') || '';
        
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//localhost:8081/api/ws/chat/${renovationId}?token=${rawToken}`;
        
        ws.current = new WebSocket(wsUrl);

        ws.current.onopen = () => {
            if (isMounted) {
                setIsConnected(true);
                setConnectionError(null);
            }
        };

        ws.current.onmessage = (event) => {
            const data = JSON.parse(event.data);
            
            if (data.error) {
                console.error("Błąd z serwera:", data.error);
                return;
            }

            if (data.status === 'sent' && data.message) {
                if (isMounted) {
                    setMessages(prev => [...prev, data.message]);
                }
                return;
            }

            if (data.ID && isMounted) {
                setMessages(prev => [...prev, data]);
            }
        };

        ws.current.onclose = (event) => {
            if (isMounted) {
                setIsConnected(false);
                if (!event.wasClean) {
                    setConnectionError("Połączenie przerwane");
                }
            }
        };

        ws.current.onerror = () => {
            if (isMounted) {
                setConnectionError("Błąd połączenia WebSocket");
                setIsConnected(false);
            }
        };

        return () => {
            isMounted = false;
            ws.current?.close();
        };
    }, [renovationId]);

    const sendMessage = (e: React.FormEvent) => {
        e.preventDefault();
        if (!inputMsg.trim() || !ws.current || !isConnected) return;

        const payload = {
            content: inputMsg.trim()
        };

        ws.current.send(JSON.stringify(payload));
        setInputMsg('');
    };

    const isMe = (senderId: string) => senderId === user?.userID;

    return (
        <div className="bg-bg-surface border border-bg-border rounded-card flex flex-col h-[600px]">
            <div className="p-4 border-b border-bg-border flex justify-between items-center bg-bg-base/30 rounded-t-card">
                <h2 className="font-bold text-text-main">Komunikator</h2>
                <div className="flex items-center gap-2 text-xs font-medium">
                    <div className={`w-2 h-2 rounded-full ${
                        isConnected ? 'bg-emerald-500' : connectionError ? 'bg-red-500' : 'bg-yellow-500'
                    }`}></div>
                    {isConnected ? 'Połączono' : connectionError ? 'Błąd połączenia' : 'Łączenie...'}
                </div>
            </div>

            <div className="flex-1 overflow-y-auto p-4 space-y-4">
                {messages.length === 0 && (
                    <div className="text-center text-text-muted text-sm mt-10">
                        Brak wiadomości. Rozpocznij konwersację.
                    </div>
                )}
                
                {messages.map((msg) => {
                    const me = isMe(msg.SenderID);
                    
                    return (
                        <div key={msg.ID || `${msg.SenderID}-${msg.CreatedAt}`} 
                             className={`flex ${me ? 'justify-end' : 'justify-start'}`}>
                            <div className={`max-w-[70%] p-3 rounded-lg text-sm ${
                                me 
                                    ? 'bg-primary text-white rounded-tr-none' 
                                    : 'bg-bg-base border border-bg-border text-text-main rounded-tl-none'
                            }`}>
                                <div>{msg.Content}</div>
                                <div className={`text-[10px] mt-1 ${me ? 'text-blue-200' : 'text-text-muted'}`}>
                                    {new Date(msg.CreatedAt).toLocaleTimeString('pl-PL', {
                                        hour: '2-digit',
                                        minute: '2-digit'
                                    })}
                                </div>
                            </div>
                        </div>
                    );
                })}
                <div ref={messagesEndRef} />
            </div>

            <form onSubmit={sendMessage} className="p-4 border-t border-bg-border bg-bg-base/30 rounded-b-card flex gap-2">
                <input
                    type="text"
                    value={inputMsg}
                    onChange={(e) => setInputMsg(e.target.value)}
                    placeholder="Wpisz wiadomość..."
                    disabled={!isConnected}
                    className="flex-1 bg-bg-surface border border-bg-border text-text-main text-sm p-2.5 rounded focus:outline-none focus:border-primary disabled:opacity-50"
                />
                <button 
                    type="submit" 
                    disabled={!isConnected || !inputMsg.trim()}
                    className="bg-primary hover:bg-primary-hover text-white px-4 py-2.5 rounded text-sm font-bold transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                    Wyślij
                </button>
            </form>
        </div>
    );
};