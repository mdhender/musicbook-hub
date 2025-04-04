// Copyright (c) 2025 Michael D Henderson. All rights reserved.

export async function getBook(id, token) {
    const headers = token
        ? { Authorization: `Bearer ${token}` }
        : {};

    const res = await fetch(`/api/books/${id}`, { headers });
    if (!res.ok) throw new Error("Failed to fetch book");
    return await res.json();
}
