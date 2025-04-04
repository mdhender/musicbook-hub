// src/App.jsx
import {useEffect, useState} from "react";
import {BrowserRouter as Router, Link, Route, Routes} from "react-router-dom";
import {API_URL, LOGIN_URL} from "./config.js";
import LoginForm from "./components/LoginForm";
import BookList from "./components/BookList";
import BookDetail from "./components/BookDetail";
import AddBookForm from "./components/AddBookForm";
import EditBookForm from "./components/EditBookForm";


function App() {
    const [books, setBooks] = useState([]);
    const [token, setToken] = useState(localStorage.getItem("token") || null);
    const [form, setForm] = useState({
        title: "",
        author: "",
        instrument: "",
        condition: "",
        description: "",
        public: false, // default to private
    });
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

    function handleInput(e) {
        const {name, type, value, checked} = e.target;
        setForm(prev => ({
            ...prev,
            [name]: type === "checkbox" ? checked : value,
        }));
    }

    function addBook(e) {
        e.preventDefault();
        fetch(API_URL, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify(form),
        })
            .then(res => {
                if (res.status === 401) throw new Error("Not authorized");
                return res.json();
            })
            .then(book => {
                setBooks(prev => [...prev, book]);
                setForm({
                    title: "",
                    author: "",
                    instrument: "",
                    condition: "",
                    public: true
                });
            })
            .catch(err => alert("Failed to add book: " + err.message));
    }

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
                                    <LoginForm token={token} login={login} logout={logout} />

                                    <AddBookForm
                                        form={form}
                                        handleInput={handleInput}
                                        addBook={addBook}
                                    />

                                    {/* Status Messages */}
                                    <div className="space-y-2">
                                        {loading && (
                                            <p className="text-gray-500 italic">Loading books...</p>
                                        )}
                                        {error && (
                                            <p className="text-red-600">
                                                Error loading books. Please try again later.
                                            </p>
                                        )}
                                        {/* Book Count */}
                                        {!loading && !error && books.length !== 0 && (
                                            <p className="text-sm text-gray-500 mb-2">{books.length} book(s) visible</p>
                                        )}
                                        {!loading && !error && books.length === 0 && (
                                            <p className="text-gray-500 italic">No books available.</p>
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
                                </>
                            }
                        />
                        <Route path="/books/:id" element={<BookDetail token={token} />} />
                        <Route path="/books/:id/edit" element={<EditBookForm token={token} />} />
                    </Routes>
                </div>
            </div>
        </Router>
    );
}

export default App;
