const API_URL = "http://localhost:8080";
let token = localStorage.getItem("token");
let user = JSON.parse(localStorage.getItem("user"));

function updateUI() {
    if (token) {
        document.getElementById("auth-section").classList.add("hidden");
        document.getElementById("dashboard-section").classList.remove("hidden");
        document.getElementById("user-display").innerText = user ? user.username : "User";
        loadProducts();
        loadOrders();
    } else {
        document.getElementById("auth-section").classList.remove("hidden");
        document.getElementById("dashboard-section").classList.add("hidden");
    }
}

async function register() {
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;
    
    try {
        const res = await fetch(`${API_URL}/auth/register`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, password })
        });
        
        if (res.ok) {
            document.getElementById("auth-message").innerText = "Registered! Please login.";
        } else {
            document.getElementById("auth-message").innerText = "Registration failed.";
        }
    } catch (e) {
        document.getElementById("auth-message").innerText = "Error: " + e.message;
    }
}

async function login() {
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;
    
    try {
        const res = await fetch(`${API_URL}/auth/login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, password })
        });
        
        if (res.ok) {
            const data = await res.json();
            token = data.token;
            user = data.user;
            localStorage.setItem("token", token);
            localStorage.setItem("user", JSON.stringify(user));
            updateUI();
        } else {
            document.getElementById("auth-message").innerText = "Login failed.";
        }
    } catch (e) {
        document.getElementById("auth-message").innerText = "Error: " + e.message;
    }
}

function logout() {
    token = null;
    user = null;
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    updateUI();
}

async function loadProducts() {
    const list = document.getElementById("products-list");
    list.innerHTML = "Loading...";
    
    try {
        const res = await fetch(`${API_URL}/products`);
        const products = await res.json();
        
        list.innerHTML = "";
        products.forEach(p => {
            const div = document.createElement("div");
            div.className = "card";
            div.innerHTML = `
                <h4>${p.name} ($${p.price})</h4>
                <p>${p.description}</p>
                <p>Stock: ${p.stock}</p>
                <button onclick="addToCart('${p.id}')">Buy 1</button>
            `;
            list.appendChild(div);
        });
    } catch (e) {
        list.innerHTML = "Error loading products.";
    }
}

async function addToCart(productId) {
    if (!token) return alert("Please login first");
    
    try {
        const res = await fetch(`${API_URL}/orders`, {
            method: "POST",
            headers: { 
                "Content-Type": "application/json",
                "Authorization": token
            },
            body: JSON.stringify([{ productId, quantity: 1 }])
        });
        
        if (res.ok) {
            alert("Order placed!");
            loadOrders();
            loadProducts(); // Update stock display
        } else {
            alert("Failed to order: " + (await res.text()));
        }
    } catch (e) {
        alert("Error: " + e.message);
    }
}

async function loadOrders() {
    const list = document.getElementById("orders-list");
    list.innerHTML = "Loading...";
    
    try {
        const res = await fetch(`${API_URL}/orders`, {
            headers: { "Authorization": token }
        });
        const orders = await res.json();
        
        list.innerHTML = "";
        if (!orders || orders.length === 0) {
            list.innerHTML = "No orders yet.";
            return;
        }

        orders.forEach(o => {
            const div = document.createElement("div");
            div.className = "card";
            div.innerHTML = `
                <h4>Order #${o.id}</h4>
                <p>Status: ${o.status}</p>
                <p>Items: ${o.items.length}</p>
                <p>Date: ${new Date(o.createdAt).toLocaleString()}</p>
            `;
            list.appendChild(div);
        });
    } catch (e) {
        list.innerHTML = "Error loading orders.";
    }
}

// Initialize
updateUI();
