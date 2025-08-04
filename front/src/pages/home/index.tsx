"use client"

import { useEffect, useState } from "react"
import { BookCard } from "@/components/bookcard"
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area"

export default function HomePage() {
	const [readingBooks, setReadingBooks] = useState([])
	const [recentlyAddedBooks, setRecentlyAddedBooks] = useState([])

	useEffect(() => {
		async function fetchBooks() {
			try {
				const res1 = await fetch(`/api/all?sort=last_opened&order=desc&page=1`)
				if (!res1.ok) throw new Error("Failed to fetch books")
				const data1 = await res1.json()
				setReadingBooks(data1.books)

				const res2 = await fetch(`/api/all?sort=added_time&order=desc&page=1`)
				if (!res2.ok) throw new Error("Failed to fetch books")
				const data2 = await res2.json()
				setRecentlyAddedBooks(data2.books)
			} catch (e) {
			} finally {
			}
		}
		fetchBooks()
	}, [])

	return (
		<div className="p-6 space-y-6">
			<h1 className="text-2xl font-semibold">Home Page</h1>

			<section className="space-y-2">
				<h2 className="text-xl font-medium">Reading</h2>
				<ScrollArea className="w-[calc(100vw-19rem)] max-w-full">
					<div className="inline-flex gap-4 ">
						{readingBooks.map((book, index) => (
							<BookCard key={index} book={book} />
						))}
					</div>
					<ScrollBar orientation="horizontal" />
				</ScrollArea>
			</section>

			<section className="space-y-2">
				<h2 className="text-xl font-medium">Arrivals</h2>
				<ScrollArea className="w-[calc(100vw-19rem)] max-w-full">
					<div className="inline-flex gap-4 ">
						{recentlyAddedBooks.map((book, index) => (
							<BookCard key={index} book={book} />
						))}
					</div>
					<ScrollBar orientation="horizontal" />
				</ScrollArea>
			</section>
		</div>
	)
}
