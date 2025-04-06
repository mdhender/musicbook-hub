// src/components/LoginForm.jsx
// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import Button from "./Button";

export default function LoginForm({ token, login, logout }) {
    return (
        <div className="mb-4">
            {!token ? (
                <form
                    onSubmit={(e) => {
                        e.preventDefault();
                        const uuid = e.target.uuid.value;
                        login(uuid);
                    }}
                    className="flex flex-col sm:flex-row items-start sm:items-center gap-2"
                >
                    <input
                        name="uuid"
                        placeholder="Enter credentials"
                        className="p-2 border w-full sm:w-2/3 rounded"
                    />
                    <Button variant="blue" type="submit">Login</Button>
                </form>
            ) : (
                <div className="flex justify-between items-center">
                    <span className="text-green-600 font-medium">✅ Logged in</span>
                    <Button variant="red" onClick={logout}>Logout</Button>
                </div>
            )}
        </div>
    );
}
