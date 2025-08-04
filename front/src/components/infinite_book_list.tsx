"use client"

import { useEffect, useState } from "react"
import { useInView } from "react-intersection-observer"
import { ScrollArea } from "@/components/ui/scroll-area"
import { BookCard } from "@/components/bookcard"
import { BookEntry, SortKey, SortOrder } from "@/api/interface"

interface InfiniteBookListProps {
    apiEndpoint: string
    sortKey: SortKey;
    sortOrder: SortOrder;
    q?: string;
}

export function InfiniteBookList({ apiEndpoint, sortKey, sortOrder, q }: InfiniteBookListProps) {
    const [books, setBooks] = useState<BookEntry[]>([])
    const [page, setPage] = useState(1)
    const [hasMore, setHasMore] = useState(true)
    const [isLoading, setIsLoading] = useState(false)

    const { ref, inView } = useInView()

    useEffect(() => {
        setBooks([])
        setPage(1)
        setHasMore(true)
    }, [apiEndpoint, sortKey, sortOrder, q])


    useEffect(() => {
        if (!hasMore || isLoading) return

        const loadBooks = async () => {
            setIsLoading(true)
            try {
                const res = await fetch(`${apiEndpoint}?sort=${sortKey}&order=${sortOrder}&page=${page}&q=${q ?? ""}`)
                if (!res.ok) throw new Error("Failed to fetch books")

                const data = await res.json()

                setBooks((prev) => [...prev, ...data.books])
                setHasMore(data.hasMore)

                setPage((p) => p + 1)
            } catch (e) {
                console.error(e)
                setHasMore(false)
            } finally {
                setIsLoading(false)
            }
        }

        loadBooks()
    }, [inView, apiEndpoint, sortKey, sortOrder, q])

    return (
        <ScrollArea className="h-[89vh] pr-2">
            <div className="flex flex-wrap gap-4">
                {books.map((book) => (
                    <BookCard key={book.path} book={book} />
                ))}

                {hasMore && (
                    <div ref={ref} className="w-full text-center py-4 text-muted-foreground">
                        {isLoading ? "Loading..." : "Show more"}
                    </div>
                )}
            </div>
        </ScrollArea>
    )
}
