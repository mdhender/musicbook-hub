// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { API_URL } from "./config.js";


export default function BookDetail({ token }) {
    const { id } = useParams();
    const [book, setBook] = useState(null);
    const [error, setError] = useState(null);

    useEffect(() => {
        const headers = token
            ? { Authorization: `Bearer ${token}` }
            : {};

        fetch(`${API_URL}/${id}`, { headers })
            .then(res => {
                if (!res.ok) throw new Error("Book not found or not authorized");
                return res.json();
            })
            .then(data => setBook(data))
            .catch(err => setError(err.message));
    }, [id, token]);

    if (error) return <p className="text-red-600">{error}</p>;
    if (!book) return <p>Loading...</p>;

    return (
        <div className="space-y-4">
            <Link to="/" className="text-blue-600 underline">&larr; Back</Link>
            <h2 className="text-2xl font-bold">{book.title}</h2>
            <div><strong>Author:</strong> {book.author}</div>
            <div><strong>Instrument:</strong> {book.instrument}</div>
            <div><strong>Condition:</strong> {book.condition}</div>
            <div><strong>Description:</strong> {book.description}</div>
            <div>
                <strong>Visibility:</strong>{" "}
                {book.public ? (
                    <span className="text-green-600 font-semibold">🌍 Public</span>
                ) : (
                    <span className="text-red-600 font-semibold">🔒 Private</span>
                )}
            </div>
        </div>
    );
}
