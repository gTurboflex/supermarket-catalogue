const API_BASE = window.location.origin;

function showTab(tabId) {
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.style.display = 'none';
    });
    
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    
    document.getElementById(tabId).style.display = 'block';
    event.currentTarget.classList.add('active');
}

function showError(elementId, error) {
    const element = document.getElementById(elementId);
    element.innerHTML = `<div class="error">Error: ${error.message || error}</div>`;
}

function showLoading(elementId) {
    const element = document.getElementById(elementId);
    element.innerHTML = `<div class="loading">Loading...</div>`;
}

function showSuccess(elementId, message) {
    const element = document.getElementById(elementId);
    element.innerHTML = `<div class="success">${message}</div>`;
}

async function makeRequest(method, endpoint, data = null) {
    const options = {
        method: method,
        headers: {
            'Content-Type': 'application/json',
        },
    };
    
    if (data) {
        options.body = JSON.stringify(data);
    }
    
    try {
        const response = await fetch(API_BASE + endpoint, options);
        
        const responseText = await response.text();
        let responseData;
        
        try {
            responseData = JSON.parse(responseText);
        } catch {
            responseData = responseText;
        }
        
        if (!response.ok) {
            throw new Error(responseData.error || responseData.message || `HTTP ${response.status}`);
        }
        
        return responseData;
    } catch (error) {
        throw error;
    }
}

function getAllProducts() {
    showLoading('productsResult');
    
    makeRequest('GET', '/products')
        .then(products => {
            if (products.length === 0) {
                document.getElementById('productsResult').innerHTML = '<p>No products found</p>';
                return;
            }
            
            let html = '<table>';
            html += '<tr><th>ID</th><th>Name</th><th>Price</th><th>Stock</th><th>Barcode</th><th>Supermarket</th><th>Actions</th></tr>';
            
            products.forEach(product => {
                html += `
                <tr>
                    <td>${product.id}</td>
                    <td>${product.name}</td>
                    <td>₸${product.price.toFixed(2)}</td>
                    <td>${product.stock}</td>
                    <td>${product.barcode || '-'}</td>
                    <td>${product.supermarket_id || '-'}</td>
                    <td>
                        <button onclick="viewProduct(${product.id})" style="padding: 5px 10px; margin: 2px; background: #17a2b8;">View</button>
                        <button onclick="deleteProduct(${product.id})" style="padding: 5px 10px; margin: 2px; background: #dc3545;">Delete</button>
                    </td>
                </tr>`;
            });
            
            html += '</table>';
            html += `<p>Total products: ${products.length}</p>`;
            
            document.getElementById('productsResult').innerHTML = html;
        })
        .catch(error => showError('productsResult', error));
}

function viewProduct(id) {
    showLoading('productsResult');
    
    makeRequest('GET', `/products/${id}`)
        .then(product => {
            const html = `
                <h3>Product Details</h3>
                <div class="json-output">${JSON.stringify(product, null, 2)}</div>
                <button onclick="getAllProducts()" style="margin-top: 15px;">Back to List</button>
            `;
            document.getElementById('productsResult').innerHTML = html;
        })
        .catch(error => showError('productsResult', error));
}

function deleteProduct(id) {
    if (!confirm(`Are you sure you want to delete product ${id}?`)) {
        return;
    }
    
    showLoading('productsResult');
    
    makeRequest('DELETE', `/products/${id}`)
        .then(() => {
            showSuccess('productsResult', `Product ${id} deleted successfully`);
            setTimeout(getAllProducts, 1500);
        })
        .catch(error => showError('productsResult', error));
}

function showCreateForm() {
    document.getElementById('createProductForm').style.display = 'block';
}

function hideCreateForm() {
    document.getElementById('createProductForm').style.display = 'none';
}

function createProduct() {
    const product = {
        name: document.getElementById('prodName').value,
        price: parseFloat(document.getElementById('prodPrice').value),
        stock: parseInt(document.getElementById('prodStock').value),
        barcode: document.getElementById('prodBarcode').value,
        image: document.getElementById('prodImage').value,
        category_id: parseInt(document.getElementById('prodCategory').value) || 1,
        supermarket_id: parseInt(document.getElementById('prodSupermarket').value) || 1,
        unit_price: parseFloat(document.getElementById('prodPrice').value)
    };
    
    if (!product.name || !product.price) {
        alert('Product name and price are required');
        return;
    }
    
    showLoading('productsResult');
    
    makeRequest('POST', '/products', product)
        .then(createdProduct => {
            showSuccess('productsResult', `Product created successfully with ID: ${createdProduct.id}`);
            hideCreateForm();
            clearForm();
            setTimeout(getAllProducts, 1500);
        })
        .catch(error => showError('productsResult', error));
}

function clearForm() {
    document.getElementById('prodName').value = '';
    document.getElementById('prodPrice').value = '';
    document.getElementById('prodStock').value = '';
    document.getElementById('prodBarcode').value = '';
    document.getElementById('prodImage').value = '';
    document.getElementById('prodCategory').value = '';
    document.getElementById('prodSupermarket').value = '';
}

function compareBarcode() {
    const barcode = document.getElementById('barcodeInput').value.trim();
    
    if (!barcode) {
        alert('Please enter a barcode');
        return;
    }
    
    showLoading('compareResult');
    
    makeRequest('GET', `/products/compare/${barcode}`)
        .then(data => {
            if (data.results.length === 0) {
                document.getElementById('compareResult').innerHTML = '<p>No products found with this barcode</p>';
                return;
            }
            
            let html = `<h3>Barcode: ${data.barcode}</h3>`;
            html += '<table>';
            html += '<tr><th>Product ID</th><th>Name</th><th>Price</th><th>Unit Price</th><th>Supermarket</th><th>Last Updated</th></tr>';
            
            data.results.forEach(item => {
                const isBest = data.best && item.product_id === data.best.product_id;
                const rowClass = isBest ? 'best-offer' : '';
                html += `
                <tr class="${rowClass}">
                    <td>${item.product_id}</td>
                    <td>${item.name}</td>
                    <td>₸${item.price.toFixed(2)}</td>
                    <td>${item.unit_price ? '₸' + item.unit_price.toFixed(2) : '-'}</td>
                    <td>${item.supermarket_name || '-'}</td>
                    <td>${item.last_updated || '-'}</td>
                </tr>`;
            });
            
            html += '</table>';
            
            if (data.best) {
                html += `
                <div class="success" style="margin-top: 20px;">
                    <strong>Best Offer:</strong> ${data.best.name} at ₸${data.best.unit_price ? data.best.unit_price.toFixed(2) : data.best.price.toFixed(2)} 
                    ${data.best.unit_price ? '(unit price)' : ''} from ${data.best.supermarket_name || 'Unknown'}
                </div>`;
            }
            
            document.getElementById('compareResult').innerHTML = html;
        })
        .catch(error => showError('compareResult', error));
}

function compareBasket() {
    const basketText = document.getElementById('basketItems').value;
    
    if (!basketText) {
        alert('Please enter basket items in JSON format');
        return;
    }
    
    let items;
    try {
        items = JSON.parse(basketText);
    } catch (e) {
        alert('Invalid JSON format. Please check your input.');
        return;
    }
    
    showLoading('basketResult');
    
    makeRequest('POST', '/basket/compare', { items: items })
        .then(data => {
            if (!data.results || data.results.length === 0) {
                document.getElementById('basketResult').innerHTML = '<p>No supermarket data found</p>';
                return;
            }
            
            let html = '<h3>Basket Comparison Results</h3>';
            html += '<table>';
            html += '<tr><th>Supermarket</th><th>Total Cost</th><th>Matched Items</th><th>Missing Items</th></tr>';
            
            data.results.forEach(supermarket => {
                const missingText = supermarket.missing.length > 0 
                    ? supermarket.missing.join(', ')
                    : 'None';
                
                html += `
                <tr>
                    <td><strong>${supermarket.supermarket_name}</strong> (ID: ${supermarket.supermarket_id})</td>
                    <td>₸${supermarket.total.toFixed(2)}</td>
                    <td>${supermarket.matched_items}</td>
                    <td>${missingText}</td>
                </tr>`;
            });
            
            html += '</table>';
            
            const bestOption = data.results.reduce((best, current) => 
                current.total < best.total ? current : best
            );
            
            html += `
            <div class="success" style="margin-top: 20px;">
                <strong>Best Option:</strong> ${bestOption.supermarket_name} - Total: ₸${bestOption.total.toFixed(2)}
                ${bestOption.missing.length > 0 ? `(Note: ${bestOption.missing.length} items missing)` : ''}
            </div>`;
            
            document.getElementById('basketResult').innerHTML = html;
        })
        .catch(error => showError('basketResult', error));
}

function getSupermarketStats() {
    showLoading('statsResult');
    
    makeRequest('GET', '/supermarkets/stats')
        .then(stats => {
            if (stats.length === 0) {
                document.getElementById('statsResult').innerHTML = '<p>No supermarket statistics available</p>';
                return;
            }
            
            let html = '<h3>Supermarket Statistics</h3>';
            html += '<table>';
            html += '<tr><th>Supermarket</th><th>Product Count</th><th>Avg Price</th><th>Min Price</th><th>Max Price</th></tr>';
            
            stats.forEach(stat => {
                html += `
                <tr>
                    <td><strong>${stat.supermarket_name}</strong> (ID: ${stat.supermarket_id})</td>
                    <td>${stat.product_count}</td>
                    <td>₸${stat.avg_price ? stat.avg_price.toFixed(2) : '0.00'}</td>
                    <td>₸${stat.min_price ? stat.min_price.toFixed(2) : '0.00'}</td>
                    <td>₸${stat.max_price ? stat.max_price.toFixed(2) : '0.00'}</td>
                </tr>`;
            });
            
            html += '</table>';
            html += `<p>Total supermarkets: ${stats.length}</p>`;
            
            document.getElementById('statsResult').innerHTML = html;
        })
        .catch(error => showError('statsResult', error));
}

function getUsers() {
    showLoading('usersResult');
    
    makeRequest('GET', '/users')
        .then(users => {
            if (users.length === 0) {
                document.getElementById('usersResult').innerHTML = '<p>No users found</p>';
                return;
            }
            
            let html = '<h3>Team Members</h3>';
            html += '<table>';
            html += '<tr><th>ID</th><th>Name</th><th>Role</th></tr>';
            
            users.forEach(user => {
                html += `
                <tr>
                    <td>${user.id}</td>
                    <td>${user.name}</td>
                    <td>${user.role}</td>
                </tr>`;
            });
            
            html += '</table>';
            html += `<p>Total team members: ${users.length}</p>`;
            
            document.getElementById('usersResult').innerHTML = html;
        })
        .catch(error => showError('usersResult', error));
}

function healthCheck() {
    showLoading('healthResult');
    
    makeRequest('GET', '/health')
        .then(health => {
            const html = `
                <div class="success">
                    <h3>✅ API is Healthy</h3>
                    <p>Status: ${health.status}</p>
                    <p>Message: ${health.message}</p>
                    <p>Timestamp: ${new Date().toLocaleString()}</p>
                </div>
                <div class="json-output" style="margin-top: 20px;">
                    ${JSON.stringify(health, null, 2)}
                </div>
            `;
            document.getElementById('healthResult').innerHTML = html;
        })
        .catch(error => showError('healthResult', error));
}

document.addEventListener('DOMContentLoaded', function() {
    getAllProducts();
});