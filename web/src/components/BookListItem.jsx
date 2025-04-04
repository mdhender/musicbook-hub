// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import { Link } from "react-router-dom";

export default function BookListItem({ book, token, toggleVisibility, deleteBook }) {
    return (
        <li className="border rounded p-3 bg-white shadow flex justify-between items-center">
            <div>
                <Link to={`/books/${book.id}`} className="text-blue-600 hover:underline">
                    <strong>{book.title}</strong>
                </Link>
                {book.author && (
                    <div className="text-sm text-gray-600">by {book.author}</div>
                )}
                {book.description && (
                    <div className="text-sm text-gray-700 mt-1">{book.description}</div>
                )}
                {token && (
                    <div
                        className="text-sm mt-1 cursor-pointer"
                        title="Click to toggle public/private"
                        onClick={() => toggleVisibility(book)}
                    >
                        {book.public ? (
                            <span className="text-green-600 font-semibold">🌍 Public</span>
                        ) : (
                            <span className="text-red-600 font-semibold">🔒 Private</span>
                        )}
                    </div>
                )}
            </div>
            {token && (
                <button
                    onClick={() => deleteBook(book.id)}
                    className="text-red-600 hover:underline"
                >
                    Delete
                </button>
            )}
        </li>
    );
}
