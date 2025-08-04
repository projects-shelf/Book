import { useState, useEffect } from "react";
import { EPUBViewer } from "@/components/viewer/epub";

function useQuery() {
    const [query, setQuery] = useState<URLSearchParams | null>(null);

    useEffect(() => {
        setQuery(new URLSearchParams(window.location.search));
    }, []);

    return query;
}

export function EPUBViewerPage() {
    const query = useQuery();

    if (!query) return <p>Loading...</p>;

    const path = query.get("path") ?? "";
    const position = query.get("position") ?? ""

    return <EPUBViewer fileUrl={`/book/epub?path=${encodeURIComponent(path)}`} initialPage={position} />;
}