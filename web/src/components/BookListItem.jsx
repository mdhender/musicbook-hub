// src/components/BookListItem.jsx
// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import {Link} from "react-router-dom";
import Button from "./Button";
import PublicToggle from "./PublicToggle";

export default function BookListItem({book, token, toggleVisibility, deleteBook}) {
    return (
        <li className="border rounded p-3 bg-white shadow flex justify-between items-start sm:items-center gap-4">
            <div className="flex-1">
                <Link to={`/books/${book.id}`}
                      className="text-blue-600 hover:underline font-semibold">
                    {book.title}
                </Link>
                {book.author && (
                    <div className="text-sm text-gray-600">
                        by {book.author}
                    </div>
                )}
                {book.description && (
                    <div className="text-sm text-gray-700 mt-1">
                        {book.description}
                    </div>
                )}
                {token && (
                    <div
                        className="text-sm mt-1 cursor-pointer select-none"
                        title="Click to toggle public/private"
                        onClick={() => toggleVisibility(book)}
                    >
                        <PublicToggle
                            isPublic={book.public}
                            onClick={() => toggleVisibility(book)}
                            interactive
                        />
                    </div>
                )}
            </div>

            {token && (
                <Button
                    variant="red"
                    onClick={() => deleteBook(book.id)}
                    className="text-sm"
                >
                    Delete
                </Button>
            )}
        </li>
    );
}
