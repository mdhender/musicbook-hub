// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import BookListItem from "./BookListItem";

export default function BookList({ books, loading, error, token, toggleVisibility, deleteBook }) {
    if (loading) {
        return <p className="text-gray-500 italic">Loading books...</p>;
    }

    if (error) {
        return <p className="text-red-600">Error loading books. Please try again later.</p>;
    }

    if (books.length === 0) {
        return <p className="text-gray-500 italic">No books available.</p>;
    }

    return (
        <ul className="space-y-2">
            {books.map(book => (
                <BookListItem
                    key={book.id}
                    book={book}
                    token={token}
                    toggleVisibility={toggleVisibility}
                    deleteBook={deleteBook}
                />
            ))}
        </ul>
    );
}
