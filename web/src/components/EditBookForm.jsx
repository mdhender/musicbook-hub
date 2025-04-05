// src/components/EditBookForm.jsx
// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";

import { API_URL } from "../config.js";
import { formatOptions } from "../formatOptions.js";

export default function EditBookForm({ token }) {
    const { id } = useParams();
    const navigate = useNavigate();
    const [form, setForm] = useState(null);
    const [error, setError] = useState(null);

    useEffect(() => {
        fetch(`${API_URL}/${id}`, {
            headers: token ? { Authorization: `Bearer ${token}` } : {},
        })
            .then(res => {
                if (!res.ok) throw new Error("Book not found or unauthorized");
                return res.json();
            })
            .then(data => setForm(data))
            .catch(err => setError(err.message));
    }, [id, token]);

    const handleInput = e => {
        const { name, value, type, checked } = e.target;
        setForm(prev => ({
            ...prev,
            [name]: type === "checkbox" ? checked : value,
        }));
    };

    const handleSubmit = e => {
        e.preventDefault();

        fetch(`${API_URL}/${id}`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(form),
        })
            .then(res => {
                if (!res.ok) throw new Error("Failed to update book");
                navigate(`/books/${id}`);
            })
            .catch(err => alert("Update failed: " + err.message));
    };

    if (error) return <p className="text-red-600">{error}</p>;
    if (!form) return <p>Loading...</p>;

    return (
        <div className="max-w-xl mx-auto p-6 space-y-4">
            <h2 className="text-2xl font-bold">Edit Book</h2>

            <form onSubmit={handleSubmit} className="space-y-3">
                <input
                    className="w-full p-2 border rounded"
                    placeholder="Title"
                    name="title"
                    value={form.title}
                    onChange={handleInput}
                    required
                />
                <input
                    className="w-full p-2 border rounded"
                    placeholder="Author"
                    name="author"
                    value={form.author}
                    onChange={handleInput}
                />
                <input
                    className="w-full p-2 border rounded"
                    placeholder="Instrument"
                    name="instrument"
                    value={form.instrument}
                    onChange={handleInput}
                />
                <input
                    className="w-full p-2 border rounded"
                    placeholder="Condition"
                    name="condition"
                    value={form.condition}
                    onChange={handleInput}
                />
                <textarea
                    className="w-full p-2 border rounded"
                    placeholder="Description"
                    name="description"
                    value={form.description}
                    onChange={handleInput}
                />
                <select
                    name="format"
                    value={form.format || ""}
                    onChange={handleInput}
                    className="w-full p-2 border rounded"
                >
                    <option value="">Select format</option>
                    {formatOptions.map((opt) => (
                        <option key={opt.format} value={opt.format}>
                            {opt.format}
                        </option>
                    ))}
                </select>
                <label className="flex items-center space-x-2">
                    <input
                        type="checkbox"
                        name="public"
                        checked={form.public}
                        onChange={handleInput}
                    />
                    <span>Public</span>
                </label>

                <div className="flex space-x-4 pt-2">
                    <button
                        type="submit"
                        className="bg-blue-600 font-semibold px-4 py-2 rounded hover:bg-blue-700"
                    >
                        Save
                    </button>
                    <Link
                        to={`/books/${id}`}
                        className="text-gray-600 underline px-4 py-2"
                    >
                        Cancel
                    </Link>
                </div>
            </form>
        </div>
    );
}
