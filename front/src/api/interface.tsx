type BookType = "EPUB" | "PDF" | "CBZ" | "CBR" | "Folder";

export interface BookEntry {
    type: BookType;
    path: string;           // path to book or folder
    cover: string;      // path to cover img
    title: string;
    currentPosition: string; // CFI or folio
    progress: number;       // 0〜100%
}
export type SortKey = "title" | "added_time" | "last_opened" | "progress";

export type SortOrder = "asc" | "desc";