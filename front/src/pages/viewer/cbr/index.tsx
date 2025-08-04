import { useState, useEffect } from "react";
import { CBRViewer } from "@/components/viewer/cbr";

function useQuery() {
    const [query, setQuery] = useState<URLSearchParams | null>(null);

    useEffect(() => {
        setQuery(new URLSearchParams(window.location.search));
    }, []);

    return query;
}

export function CBRViewerPage() {
    const query = useQuery();

    if (!query) return <p>Loading...</p>;

    const path = query.get("path") ?? "";
    const parsed = parseInt(query.get("position") ?? "", 10);
    const position = isNaN(parsed) ? 1 : parsed;

    return <CBRViewer fileUrl={`/book/cbr?path=${encodeURIComponent(path)}`} initialPage={position} />;
}
