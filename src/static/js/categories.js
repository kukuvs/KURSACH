const API_BASE = "/api";
let categoryModal;

document.addEventListener("DOMContentLoaded", () => {
	categoryModal = new bootstrap.Modal(
		document.getElementById("categoryModal")
	);
	loadCategories();
});

function openCategoryModal() {
	document.getElementById("categoryForm").reset();
	document.getElementById("categoryId").value = "";
	document.getElementById("categoryModalTitle").textContent =
		"Новая категория";
	categoryModal.show();
}

async function loadCategories() {
	try {
		const response = await fetch(`${API_BASE}/categories`);
		const categories = await response.json();
		const tbody = document.getElementById("categoriesTableBody");
		tbody.innerHTML = "";

		if (categories && categories.length > 0) {
			categories.forEach((cat) => {
				const row = document.createElement("tr");
				row.innerHTML = `
                    <td><strong>${cat.name}</strong></td>
                    <td>${cat.description || ""}</td>
                    <td>
                        <span style="display:inline-block; width:20px; height:20px; background-color:${
							cat.color
						}; border-radius:50%;"></span>
                    </td>
                    <td class="text-end">
                        <button class="btn btn-sm btn-outline-primary me-1" onclick="editCategory(${
							cat.id
						})">
                            <i class="bi bi-pencil"></i>
                        </button>
                        <button class="btn btn-sm btn-outline-danger" onclick="deleteCategory(${
							cat.id
						})">
                            <i class="bi bi-trash"></i>
                        </button>
                    </td>
                `;
				tbody.appendChild(row);
			});
		} else {
			tbody.innerHTML =
				'<tr><td colspan="4" class="text-center text-muted">Нет категорий</td></tr>';
		}
	} catch (error) {
		console.error("Error loading categories:", error);
	}
}

async function saveCategory() {
	const id = document.getElementById("categoryId").value;
	const data = {
		name: document.getElementById("categoryName").value,
		description: document.getElementById("categoryDescription").value,
		color: document.getElementById("categoryColor").value,
	};

	try {
		let url = `${API_BASE}/categories`;
		let method = "POST";
		if (id) {
			url += "/" + id;
			method = "PUT";
		}

		const response = await fetch(url, {
			method: method,
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(data),
		});

		if (response.ok) {
			categoryModal.hide();
			loadCategories();
		} else {
			alert("Ошибка сохранения");
		}
	} catch (e) {
		console.error(e);
	}
}

async function editCategory(id) {
	try {
		const response = await fetch(`${API_BASE}/categories`);
		const categories = await response.json();
		const cat = categories.find((c) => c.id == id);

		if (cat) {
			document.getElementById("categoryId").value = cat.id;
			document.getElementById("categoryName").value = cat.name;
			document.getElementById("categoryDescription").value =
				cat.description;
			document.getElementById("categoryColor").value = cat.color;
			document.getElementById("categoryModalTitle").textContent =
				"Редактирование категории";
			categoryModal.show();
		}
	} catch (e) {
		console.error(e);
	}
}

async function deleteCategory(id) {
	if (!confirm("Удалить категорию?")) return;
	try {
		const response = await fetch(`${API_BASE}/categories/${id}`, {
			method: "DELETE",
		});
		if (response.ok) loadCategories();
	} catch (e) {
		console.error(e);
	}
}
