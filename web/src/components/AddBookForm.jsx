// Copyright (c) 2025 Michael D Henderson. All rights reserved.

export default function AddBookForm({ form, handleInput, addBook }) {
    return (
        <form onSubmit={addBook} className="space-y-3 mb-6">
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
                placeholder="Author (optional)"
                name="author"
                value={form.author}
                onChange={handleInput}
            />
            <input
                className="w-full p-2 border rounded"
                placeholder="Instrument (optional)"
                name="instrument"
                value={form.instrument}
                onChange={handleInput}
            />
            <input
                className="w-full p-2 border rounded"
                placeholder="Condition (optional)"
                name="condition"
                value={form.condition}
                onChange={handleInput}
            />
            <label className="flex items-center space-x-2">
                <input
                    type="checkbox"
                    name="public"
                    checked={form.public}
                    onChange={handleInput}
                />
                <span>Public</span>
            </label>
            <button
                className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
                type="submit"
            >
                Add Book
            </button>
        </form>
    );
}
