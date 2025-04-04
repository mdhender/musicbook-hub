// Copyright (c) 2025 Michael D Henderson. All rights reserved.

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
                >
                    <input
                        name="uuid"
                        placeholder="Enter magic UUID"
                        className="p-2 border mr-2 w-2/3"
                    />
                    <button
                        type="submit"
                        className="bg-blue-600 px-4 py-2 rounded"
                    >
                        Login
                    </button>
                </form>
            ) : (
                <div className="flex justify-between items-center">
                    <span className="text-green-600">✅ Logged in</span>
                    <button
                        className="text-sm text-red-600 underline ml-4"
                        onClick={logout}
                    >
                        Logout
                    </button>
                </div>
            )}
        </div>
    );
}
