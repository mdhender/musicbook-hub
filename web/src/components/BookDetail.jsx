// src/components/BookDetail.jsx
// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";

import { API_URL } from "../config.js";

export default function BookDetail({ token }) {
    const { id } = useParams();
    const navigate = useNavigate();
    const [book, setBook] = useState(null);
    const [error, setError] = useState(null);

    useEffect(() => {
        const headers = token ? { Authorization: `Bearer ${token}` } : {};

        fetch(`${API_URL}/${id}`, { headers })
            .then(res => {
                if (!res.ok) throw new Error("Book not found or not authorized");
                return res.json();
            })
            .then(setBook)
            .catch(err => setError(err.message));
    }, [id, token]);

    const handleDelete = () => {
        if (!window.confirm("Are you sure you want to delete this book?")) return;

        fetch(`${API_URL}/${id}`, {
            method: "DELETE",
            headers: {
                Authorization: `Bearer ${token}`,
            },
        })
            .then(res => {
                if (!res.ok) throw new Error("Failed to delete book");
                navigate("/"); // Go back to home
            })
            .catch(err => alert("Delete failed: " + err.message));
    };

    const handleEdit = () => navigate(`/books/${id}/edit`);

    if (error) return <p className="text-red-600">{error}</p>;
    if (!book) return <p>Loading...</p>;

    return (
        <div className="space-y-4">
            <Link to="/" className="text-blue-600 underline">&larr; Back</Link>

            <h2 className="text-2xl font-bold">{book.title}</h2>

            {book.author && <div><strong>Author:</strong> {book.author}</div>}
            {book.instrument && <div><strong>Instrument:</strong> {book.instrument}</div>}
            {book.condition && <div><strong>Condition:</strong> {book.condition}</div>}
            {book.description && <div><strong>Description:</strong> {book.description}</div>}
            {book.format && <div><strong>Format:</strong> {book.format}</div>}

            <div>
                <strong>Visibility:</strong>{" "}
                {book.public ? (
                    <span className="text-green-600 font-semibold">🌍 Public</span>
                ) : (
                    <span className="text-red-600 font-semibold">🔒 Private</span>
                )}
            </div>

            {token && (
                <>
                    {book.created_at && (
                        <div className="text-sm text-gray-500">
                            <strong>Created:</strong> {new Date(book.created_at).toLocaleString()}
                        </div>
                    )}
                    {book.updated_at && (
                        <div className="text-sm text-gray-500">
                            <strong>Updated:</strong> {new Date(book.updated_at).toLocaleString()}
                        </div>
                    )}

                    <div className="flex space-x-4 pt-4">
                        <button
                            onClick={handleEdit}
                            className="bg-yellow-500 font-semibold px-4 py-2 rounded hover:bg-yellow-600 focus:outline-none"
                        >
                            Edit
                        </button>
                        <button
                            onClick={handleDelete}
                            className="bg-red-600 px-4 py-2 rounded hover:bg-red-700"
                        >
                            Delete
                        </button>
                    </div>
                </>
            )}
        </div>
    );
}
