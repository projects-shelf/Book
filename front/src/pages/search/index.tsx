import { useState } from "react"
import { Input } from "@/components/ui/input"  // 追加
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Switch } from "@/components/ui/switch"
import { Label } from "@/components/ui/label"
import type { SortKey, SortOrder } from "@/api/interface"
import { InfiniteBookList } from "@/components/infinite_book_list"

const sortOptions: { label: string; value: SortKey }[] = [
	{ label: "Title", value: "title" },
	{ label: "Date Added", value: "added_time" },
	{ label: "Last Opened", value: "last_opened" },
	{ label: "Progress", value: "progress" },
]

export default function SearchPage() {
	const [sortKey, setSortKey] = useState<SortKey>("title")
	const [sortOrder, setSortOrder] = useState<SortOrder>("asc")
	const [searchQuery, setSearchQuery] = useState("")

	return (
		<div className="p-6 space-y-6 w-full">
			<div className="flex items-center gap-4 w-full">
				<Input
					type="search"
					placeholder="Search books..."
					value={searchQuery}
					onChange={(e) => setSearchQuery(e.target.value)}
					className="flex-grow max-w-[33rem]"
				/>

				<div className="flex items-center gap-6 ml-auto flex-shrink-0">
					<Select value={sortKey} onValueChange={(v) => setSortKey(v as SortKey)}>
						<SelectTrigger className="w-[200px]">
							<SelectValue placeholder="Sort by" />
						</SelectTrigger>
						<SelectContent>
							{sortOptions.map((opt) => (
								<SelectItem key={opt.value} value={opt.value}>
									{opt.label}
								</SelectItem>
							))}
						</SelectContent>
					</Select>

					<div className="flex items-center gap-2">
						<Switch
							id="order-switch"
							checked={sortOrder === "asc"}
							onCheckedChange={(checked) =>
								setSortOrder(checked ? "asc" : "desc")
							}
						/>
						<Label className="w-[80px]">{sortOrder === "asc" ? "Ascending" : "Descending"}</Label>
					</div>
				</div>
			</div>

			<InfiniteBookList
				apiEndpoint="/api/search"
				sortKey={sortKey}
				sortOrder={sortOrder}
				q={encodeURIComponent(searchQuery.trim())}
			/>
		</div>
	)
}
