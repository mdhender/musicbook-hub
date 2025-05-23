// src/components/AddBookForm.jsx
// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import { useRef, useState } from "react";
import Button from "./Button";
import { API_URL } from "../config.js";
import { formatOptions } from "../formatOptions.js";

export default function AddBookForm({ token, setBooks }) {
    const titleRef = useRef();

    const [form, setForm] = useState({
        title: "",
        author: "",
        instrument: "",
        condition: "",
        description: "",
        format: "",
        public: false,
    });

    const handleInput = (e) => {
        const { name, type, value, checked } = e.target;
        setForm((prev) => ({
            ...prev,
            [name]: type === "checkbox" ? checked : value,
        }));
    };

    const handleSubmit = (e) => {
        e.preventDefault();

        fetch(API_URL, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(form),
        })
            .then((res) => {
                if (res.status === 401) throw new Error("Not authorized");
                return res.json();
            })
            .then((newBook) => {
                setBooks((prev) => [...prev, newBook]);

                const lastInstrument = form.instrument;
                setForm({
                    title: "",
                    author: "",
                    instrument: lastInstrument,
                    condition: "",
                    description: "",
                    format: "",
                    public: false,
                });

                titleRef.current?.focus();
            })
            .catch((err) => alert("Failed to add book: " + err.message));
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-3 mb-6">
            <input
                ref={titleRef}
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
                value={form.format}
                onChange={handleInput}
                className="w-full p-2 border rounded"
                required
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
            <Button type="submit" variant="blue">Add Book</Button>
        </form>
    );
}
