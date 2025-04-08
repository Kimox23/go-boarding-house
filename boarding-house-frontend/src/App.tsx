import { createSignal, createEffect, For } from 'solid-js';

function App() {
  const [users, setUsers] = createSignal(null); // Reactive state for users
  const [error, setError] = createSignal(null); // Reactive state for error

  const fetchData = async () => {
    try {
      const response = await fetch('http://localhost:8000/api/users', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlIjoiYWRtaW4iLCJpc3MiOiJib2FyZGluZy1ob3VzZS1zeXN0ZW0iLCJleHAiOjE3NDQxNzAxOTksImlhdCI6MTc0NDA4Mzc5OX0.tksIwkPu9g7-bPy9AyMxY-jgkmFKGKWoLEB487s-OFw'
        }
      });

      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

      const data = await response.json();
      setUsers(data.data);
    } catch (err:any) {
      setError(err.message);
    }
  };
  // Use createEffect to fetch data on component mount
  createEffect(() => {
    fetchData();
  });

  return (
    <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100 gap-4">
      <h1 class="text-3xl font-bold underline">All Users</h1>
      {error() && <div class="text-red-500">{error()}</div>} 
      <For each={users()} fallback={<div>Loading...</div>}>
        {(user) => (
          <div class="bg-white shadow-md rounded-lg p-4 mb-4">
            <h2 class="text-xl font-semibold">Username: {user.username}</h2>
            <p class="text-gray-600">Email: {user.email}</p>
            <p class="text-gray-600">Phone: {user.phone}</p>
            <p class="text-gray-600">Role: {user.role}</p>
          </div>
        )}
      </For>
    </div>
  );
}

export default App;