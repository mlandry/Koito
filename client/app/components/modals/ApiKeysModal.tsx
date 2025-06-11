import { useQuery } from "@tanstack/react-query";
import { createApiKey, deleteApiKey, getApiKeys, type ApiKey } from "api/api";
import { AsyncButton } from "../AsyncButton";
import { useEffect, useState } from "react";
import { Copy, Trash } from "lucide-react";

type CopiedState = {
    x: number;
    y: number;
    visible: boolean;
};

export default function ApiKeysModal() {
    const [input, setInput] = useState('')
    const [loading, setLoading ] = useState(false)
    const [err, setError ] = useState<string>()
    const [displayData, setDisplayData] = useState<ApiKey[]>([])
    const [copied, setCopied] = useState<CopiedState | null>(null);
        
    const { isPending, isError, data, error } = useQuery({ 
        queryKey: [
            'api-keys'
        ], 
        queryFn: () => {
            return getApiKeys();
        },
    });

    useEffect(() => {
        if (data) {
            setDisplayData(data)
        }
    }, [data])    

    if (isError) {
        return (
            <p className="error">Error: {error.message}</p>
        )
    }
    if (isPending) {
        return (
            <p>Loading...</p>
        )
    }

    const handleCopy = (e: React.MouseEvent<HTMLButtonElement>, text: string) => {
        navigator.clipboard.writeText(text);
    
        const parentRect = (e.currentTarget.closest(".relative") as HTMLElement).getBoundingClientRect();
        const buttonRect = e.currentTarget.getBoundingClientRect();
    
        setCopied({
            x: buttonRect.left - parentRect.left + buttonRect.width / 2, // center of button
            y: buttonRect.top - parentRect.top - 8, // above the button
            visible: true,
        });
    
        setTimeout(() => setCopied(null), 1500);
    };    
    
    const handleCreateApiKey = () => {
        setError(undefined)
        if (input === "") {
            setError("a label must be provided")
            return
        }
        setLoading(true)
        createApiKey(input)
        .then(r => {
            setDisplayData([r, ...displayData])
            setInput('')
        }).catch((err) => setError(err.message))
        setLoading(false)
    }

    const handleDeleteApiKey = (id: number) => {
        setError(undefined)
        setLoading(true)
        deleteApiKey(id)
        .then(r => {
            if (r.ok) {
                setDisplayData(displayData.filter((v) => v.id != id))
            } else {
                r.json().then((r) => setError(r.error))
            }
        })
        setLoading(false)

    }

    return (
        <div className="">
        <h2>API Keys</h2>
        <div className="flex flex-col gap-4 relative">
            {displayData.map((v) => (
                <div className="flex gap-2">
                    <div className="bg p-3 rounded-md flex-grow" key={v.key}>{v.key.slice(0, 8)+'...'} {v.label}</div>
                    <button onClick={(e) => handleCopy(e, v.key)} className="large-button px-5 rounded-md"><Copy size={16} /></button>
                    <AsyncButton loading={loading} onClick={() => handleDeleteApiKey(v.id)} confirm><Trash size={16} /></AsyncButton>
                </div>
            ))}
            <div className="flex gap-2 w-3/5">
                <input
                    type="text"
                    placeholder="Add a label for a new API key"
                    className="mx-auto fg bg rounded-md p-3 flex-grow"
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                />
                <AsyncButton loading={loading} onClick={handleCreateApiKey}>Create</AsyncButton>
            </div>
            {err && <p className="error">{err}</p>}
            {copied?.visible && (
                <div
                    style={{
                        position: "absolute",
                        top: copied.y,
                        left: copied.x,
                        transform: "translate(-50%, -100%)",
                    }}
                    className="pointer-events-none bg-black text-white text-sm px-2 py-1 rounded shadow-lg opacity-90 animate-fade"
                >
                    Copied!
                </div>
            )}
        </div>
        </div>
    )
}