import { useState, useEffect } from "react"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Switch } from "@/components/ui/switch"
import {
    Select,
    SelectTrigger,
    SelectContent,
    SelectItem,
    SelectValue,
} from "@/components/ui/select"
import {
    SheetContent,
    SheetHeader,
    SheetTitle,
} from "@/components/ui/sheet"

export type ViewerOptions = {
    direction: "ltr" | "rtl"
    spread: "none" | "odd" | "even"
    fontSize: number
}

type ViewerOptionSheetProps = {
    onOptionChanged?: (options: ViewerOptions) => void
}

export function ViewerOptionSheet({ onOptionChanged }: ViewerOptionSheetProps) {
    const [direction, setDirection] = useState<"ltr" | "rtl">("ltr")
    const [spread, setSpread] = useState<"none" | "odd" | "even">("none")
    const [fontSize, setFontSize] = useState<number | "">(16)

    useEffect(() => {
        const saved = localStorage.getItem("viewerOptions")
        if (saved) {
            try {
                const parsed = JSON.parse(saved) as ViewerOptions
                setDirection(parsed.direction)
                setSpread(parsed.spread)
                setFontSize(parsed.fontSize)
            } catch (e) {
                console.warn("Failed to parse viewerOptions from localStorage")
            }
        }
    }, [])

    const updateOptions = (next: Partial<ViewerOptions>) => {
        const updated: ViewerOptions = {
            direction,
            spread,
            fontSize: fontSize === "" ? 8 : fontSize,
            ...next,
        }
        onOptionChanged?.(updated)
        localStorage.setItem("viewerOptions", JSON.stringify(updated))
    }

    return (
        <SheetContent>
            <SheetHeader>
                <SheetTitle>Options</SheetTitle>
            </SheetHeader>

            <div className="grid flex-1 auto-rows-min gap-6 px-4 pt-6">
                <div className="flex items-center justify-between">
                    <Label htmlFor="direction-switch">Right to Left</Label>
                    <Switch
                        id="direction-switch"
                        checked={direction === "rtl"}
                        onCheckedChange={(v) => {
                            const newDirection = v ? "rtl" : "ltr"
                            setDirection(newDirection)
                            updateOptions({ direction: newDirection })
                        }}
                    />
                </div>

                <div className="grid gap-3">
                    <Label htmlFor="spread-select">Spread Mode (Not supported in ePUB)</Label>
                    <Select
                        value={spread}
                        onValueChange={(v) => {
                            const newSpread = v as "none" | "odd" | "even"
                            setSpread(newSpread)
                            updateOptions({ spread: newSpread })
                        }}
                    >
                        <SelectTrigger id="spread-select">
                            <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="none">No Spreads</SelectItem>
                            <SelectItem value="odd">Odd Spreads</SelectItem>
                            <SelectItem value="even">Even Spreads</SelectItem>
                        </SelectContent>
                    </Select>
                </div>

                <div className="grid gap-3">
                    <Label htmlFor="font-size">Font Size (Only for ePUB)</Label>
                    <Input
                        id="font-size"
                        type="number"
                        min={8}
                        max={48}
                        value={fontSize}
                        onChange={(e) => {
                            if (e.target.value === "") {
                                setFontSize("")
                                return
                            }
                            const newSize = parseInt(e.target.value, 10)
                            if (!isNaN(newSize)) {
                                setFontSize(newSize)
                            }
                        }}
                        onBlur={() => {
                            if (typeof fontSize === "number") {
                                updateOptions({ fontSize })
                            }
                        }}
                        onKeyDown={(e) => {
                            if (e.key === "Enter") {
                                e.currentTarget.blur()
                            }
                        }}
                    />
                </div>
            </div>
        </SheetContent>
    )
}
