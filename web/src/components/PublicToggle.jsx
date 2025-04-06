// src/components/PublicToggle.jsx
// Copyright (c) 2025 Michael D Henderson. All rights reserved.

export default function PublicToggle({ isPublic, onClick, interactive = false }) {
    const label = isPublic ? "🌍 Public" : "🔒 Private";
    const color = isPublic ? "text-green-600" : "text-red-600";
    const classes = `text-sm font-semibold ${color} ${interactive ? "cursor-pointer select-none" : ""}`;

    return (
        <div
            title={interactive ? "Click to toggle public/private" : ""}
            className={classes}
            onClick={interactive ? onClick : undefined}
        >
            {label}
        </div>
    );
}
