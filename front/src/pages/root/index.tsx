import { useEffect, useState } from "react";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import type { SortKey, SortOrder } from "@/api/interface";
import { InfiniteBookList } from "@/components/infinite_book_list";
import { useLocation } from "react-router-dom";

const sortOptions: { label: string; value: SortKey }[] = [
	{ label: "Title", value: "title" },
	{ label: "Date Added", value: "added_time" },
	{ label: "Last Opened", value: "last_opened" },
	{ label: "Progress", value: "progress" },
];

export default function RootPage() {
	const location = useLocation()
	const [decodedName, setDecodedName] = useState("")
	const [sortKey, setSortKey] = useState<SortKey>("title")
	const [sortOrder, setSortOrder] = useState<SortOrder>("asc")
	const [apiEndpoint, setApiEndpoint] = useState("/api/root")

	useEffect(() => {
		const path = location.pathname
		setApiEndpoint(`/api${path}`)

		if (path === "/root") {
			setDecodedName("Root")
			return
		}

		const segments = path.split("/")
		const lastSegment = segments[segments.length - 1] || segments[segments.length - 2]
		const decoded = decodeURIComponent(lastSegment)
		setDecodedName(decoded)
	}, [location.pathname])

	return (
		<div className="p-6 space-y-6 w-full">
			<div className="flex items-center gap-4 w-full">
				<h1 className="text-2xl font-semibold flex-shrink-0">{decodedName}</h1>

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
				key={`${apiEndpoint}-${sortKey}-${sortOrder}`}
				apiEndpoint={apiEndpoint}
				sortKey={sortKey}
				sortOrder={sortOrder}
			/>
		</div>
	);
}
