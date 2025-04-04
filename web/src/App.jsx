// src/App.jsx
import { useEffect, useState } from "react";

// Use environment variables to determine the base URL
const API_URL = import.meta.env.PROD
    ? "/api/books"  // Production: use relative path
    : "http://localhost:8181/api/books";  // Development: use absolute path with localhost
const LOGIN_URL = import.meta.env.PROD
    ? "/api/login"  // Production: use relative path
    : "http://localhost:8181/api/login";  // Development: use absolute path with localhost

function App() {
  const [books, setBooks] = useState([]);
  const [token, setToken] = useState(localStorage.getItem("token") || null);
  const [form, setForm] = useState({
    title: "",
    author: "",
    instrument: "",
    condition: "",
    public: true, // default to public
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);


  useEffect(() => {
    const headers = token
        ? { Authorization: `Bearer ${token}` }
        : {};

    setLoading(true);     // show loading indicator
    setError(null);       // reset any previous errors

    fetch(API_URL, { headers })
        .then(res => {
          if (!res.ok) throw new Error(`Server error: ${res.status}`);
          return res.json();
        })
        .then(data => setBooks(data.books || []))
        .catch(err => setError(err))
        .finally(() => setLoading(false));
  }, [token]);

  function handleInput(e) {
    const { name, type, value, checked } = e.target;
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
    const updated = { public: !book.public };

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
      <div className="max-w-xl mx-auto p-6">
        <h1 className="text-3xl font-bold mb-4">📚 Music Book Hub</h1>

        <div className="mb-4">
          {!token ? (
              <form onSubmit={e => {
                e.preventDefault();
                const uuid = e.target.uuid.value;
                login(uuid);
              }}>
                <input name="uuid" placeholder="Enter magic UUID" className="p-2 border mr-2 w-2/3" />
                <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded">Login</button>
              </form>
          ) : (
              <div className="flex justify-between items-center">
                <span className="text-green-600">✅ Logged in</span>
                <button className="text-sm text-red-600 underline ml-4" onClick={logout}>Logout</button>
              </div>
          )}
        </div>

        {token && (
            <form onSubmit={addBook} className="space-y-3 mb-6">
              <input className="w-full p-2 border rounded" placeholder="Title" name="title" value={form.title} onChange={handleInput} required />
              <input className="w-full p-2 border rounded" placeholder="Author" name="author" value={form.author} onChange={handleInput} />
              <input className="w-full p-2 border rounded" placeholder="Instrument" name="instrument" value={form.instrument} onChange={handleInput} />
              <input className="w-full p-2 border rounded" placeholder="Condition" name="condition" value={form.condition} onChange={handleInput} />
              <label className="flex items-center space-x-2">
                <input type="checkbox" name="public" checked={form.public} onChange={handleInput} />
                <span>Public</span>
              </label>
              <button className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700" type="submit">Add Book</button>
            </form>
        )}

        {!loading && !error && (
            <p className="text-sm text-gray-500 mb-2">{books.length} book(s) visible</p>
        )}

        <div className="space-y-2">
          {loading && (
              <p className="text-gray-500 italic">Loading books...</p>
          )}

          {error && (
              <p className="text-red-600">
                Error loading books. Please try again later.
              </p>
          )}

          {!loading && !error && books.length === 0 && (
              <p className="text-gray-500 italic">No books available.</p>
          )}

          {!loading && !error && books.map(book => (
              <li key={book.id} className="border p-3 rounded flex justify-between items-center">
                <div>
                  <strong>{book.title}</strong> – {book.instrument} ({book.condition})
                  {book.author && (
                      <div className="text-sm text-gray-600">by {book.author}</div>
                  )}
                  {token && (
                      <div
                           className="text-sm mt-1 cursor-pointer"
                           title="Click to toggle public/private"
                           onClick={() => toggleVisibility(book)}
                      >
                        {book.public ? (
                            <span className="text-green-600 font-semibold">🌍 Public</span>
                        ) : (
                            <span className="text-red-600 font-semibold">🔒 Private</span>
                        )}
                      </div>
                  )}
                </div>
                {token && (
                    <button
                        onClick={() => deleteBook(book.id)}
                        className="text-red-600 hover:underline"
                    >
                      Delete
                    </button>
                )}
              </li>
          ))}
        </div>
      </div>
  );
}

export default App;
