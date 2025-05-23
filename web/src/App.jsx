// src/App.jsx
import {useEffect, useState} from "react";
import {BrowserRouter as Router, Link, Route, Routes, Navigate} from "react-router-dom";
import {API_URL, LOGIN_URL} from "./config.js";
import LoginForm from "./components/LoginForm";
import BookList from "./components/BookList";
import BookDetail from "./components/BookDetail";
import AddBookForm from "./components/AddBookForm";
import EditBookForm from "./components/EditBookForm";


function App() {
    const [books, setBooks] = useState([]);
    const [token, setToken] = useState(localStorage.getItem("token") || null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);


    useEffect(() => {
        const headers = token
            ? {Authorization: `Bearer ${token}`}
            : {};

        setLoading(true);     // show loading indicator
        setError(null);       // reset any previous errors

        fetch(API_URL, {headers})
            .then(res => {
                if (!res.ok) throw new Error(`Server error: ${res.status}`);
                return res.json();
            })
            .then(data => setBooks(data.books || []))
            .catch(err => setError(err))
            .finally(() => setLoading(false));
    }, [token]);

    function deleteBook(id) {
        fetch(`${API_URL}/${id}`, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${token}`
            }
        }).then(res => {
            if (res.status === 401) {
                alert("Not authorized to delete");
                return;
            }
            setBooks(prev => prev.filter(book => book.id !== id));
        });
    }

    function login(uuid) {
        fetch(`${LOGIN_URL}/${uuid}`)
            .then(res => res.json())
            .then(data => {
                if (data.token) {
                    localStorage.setItem("token", data.token);
                    setToken(data.token);
                } else {
                    alert("Login failed: No token received");
                }
            });
    }

    function logout() {
        localStorage.removeItem("token");
        setToken(null);
    }

    function toggleVisibility(book) {
        const updated = {public: !book.public};

        fetch(`${API_URL}/${book.id}`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify(updated),
        })
            .then(res => {
                if (!res.ok) throw new Error("Failed to update book");
                return res.json();
            })
            .then(updatedBook => {
                setBooks(prev =>
                    prev.map(b => (b.id === updatedBook.id ? updatedBook : b))
                );
            })
            .catch(err => {
                alert("Update failed: " + err.message);
            });
    }

    return (
        <Router>
            <div className="h-screen bg-white text-gray-900 px-4 w-full">
                <div className="max-w-xl w-full p-6">
                    <h1 className="text-3xl font-bold mb-4">
                        <Link to="/">📚 Music Book Hub</Link>
                    </h1>

                    <Routes>
                        <Route
                            path="/"
                            element={
                                <>
                                    <LoginForm token={token} login={login} logout={logout}/>

                                    {/* Show AddBookForm only when authenticated */}
                                    {token ? (
                                        <>
                                            <AddBookForm token={token} setBooks={setBooks}/>
                                            <div className="mt-4 mb-4">
                                                <Link
                                                    to="/books"
                                                    className="inline-block hover:bg-blue-700 text-black font-medium py-2 px-4 rounded"
                                                >
                                                    View All Books
                                                </Link>
                                            </div>
                                        </>
                                    ) : (
                                        <div className="space-y-2">
                                            {loading && (
                                                <p className="text-gray-500 italic">Loading books...</p>
                                            )}
                                            {error && (
                                                <p className="text-red-600">
                                                    Error loading books. Please try again later.
                                                </p>
                                            )}

                                            <BookList
                                                books={books}
                                                loading={loading}
                                                error={error}
                                                token={token}
                                                toggleVisibility={toggleVisibility}
                                                deleteBook={deleteBook}
                                            />
                                        </div>
                                    )}
                                </>
                            }
                        />
                        <Route
                            path="/books"
                            element={
                                token ? (
                                    <BookList
                                        books={books}
                                        loading={loading}
                                        error={error}
                                        token={token}
                                        toggleVisibility={toggleVisibility}
                                        deleteBook={deleteBook}
                                    />
                                ) : (
                                    <Navigate to="/" replace />  // Redirect to home if not authenticated
                                )
                            }
                        />
                        <Route path="/books/:id" element={<BookDetail token={token}/>}/>
                        <Route path="/books/:id/edit" element={<EditBookForm token={token}/>}/>
                    </Routes>
                </div>
            </div>
        </Router>
    );
}

export default App;
