// API Configuration
const API_BASE = 'http://localhost:1999/api/v1';

// Global state
let currentAccessToken = '';

// DOM Elements
const statusDot = document.getElementById('statusDot');
const statusText = document.getElementById('statusText');
const loadingOverlay = document.getElementById('loadingOverlay');
const toastContainer = document.getElementById('toastContainer');

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    initializeTabs();
    initializeForms();
    checkAPIHealth();
    generateRequestId();
});

// Tab Management
function initializeTabs() {
    const tabButtons = document.querySelectorAll('.tab-btn');
    const tabPanes = document.querySelectorAll('.tab-pane');

    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const targetTab = button.getAttribute('data-tab');
            
            // Remove active class from all buttons and panes
            tabButtons.forEach(btn => btn.classList.remove('active'));
            tabPanes.forEach(pane => pane.classList.remove('active'));
            
            // Add active class to clicked button and corresponding pane
            button.classList.add('active');
            document.getElementById(targetTab).classList.add('active');
        });
    });
}

// Form Initialization
function initializeForms() {
    // Auth form
    document.getElementById('authForm').addEventListener('submit', handleAuthSubmit);
    
    // User info form
    document.getElementById('userForm').addEventListener('submit', handleUserInfoSubmit);
    
    // Cards form
    document.getElementById('cardsForm').addEventListener('submit', handleCardsSubmit);
    
    // Payment form
    document.getElementById('paymentForm').addEventListener('submit', handlePaymentSubmit);
}

// Generate unique request ID
function generateRequestId() {
    const timestamp = Date.now();
    const random = Math.floor(Math.random() * 1000);
    const requestId = `req_${timestamp}_${random}`;
    document.getElementById('requestId').value = requestId;
}

// API Health Check
async function checkAPIHealth() {
    try {
        const response = await fetch(`${API_BASE}/health`);
        const data = await response.json();
        
        if (response.ok && data.status === 'ok') {
            updateStatus('online', 'API Online');
        } else {
            updateStatus('offline', 'API Error');
        }
    } catch (error) {
        updateStatus('offline', 'API Offline');
    }
}

// Update status indicator
function updateStatus(status, text) {
    statusDot.className = `status-dot ${status === 'offline' ? 'offline' : ''}`;
    statusText.textContent = text;
}

// Show/Hide loading overlay
function showLoading(show = true) {
    loadingOverlay.classList.toggle('show', show);
}

// Show toast notification
function showToast(message, type = 'success') {
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    
    const icon = type === 'error' ? 'fas fa-exclamation-circle' : 
                 type === 'warning' ? 'fas fa-exclamation-triangle' : 
                 'fas fa-check-circle';
    
    toast.innerHTML = `
        <i class="${icon}"></i>
        <span>${message}</span>
    `;
    
    toastContainer.appendChild(toast);
    
    // Auto remove after 5 seconds
    setTimeout(() => {
        toast.remove();
    }, 5000);
}

// Display result in result section
function displayResult(elementId, data, isError = false) {
    const resultSection = document.getElementById(elementId);
    const resultClass = isError ? 'result-error' : 'result-success';
    const icon = isError ? 'fas fa-exclamation-circle' : 'fas fa-check-circle';
    const title = isError ? 'Error' : 'Success';
    
    resultSection.className = `result-section ${resultClass} show`;
    resultSection.innerHTML = `
        <div class="result-header">
            <i class="${icon}"></i>
            ${title}
        </div>
        <div class="result-data">${JSON.stringify(data, null, 2)}</div>
    `;
}

// Auto-fill access token in other forms
function autoFillAccessToken(token) {
    currentAccessToken = token;
    document.getElementById('userAccessToken').value = token;
    document.getElementById('cardsAccessToken').value = token;
    document.getElementById('paymentAccessToken').value = token;
    
    showToast('Access token auto-filled in other forms', 'success');
}

// Handle Authentication Form Submit
async function handleAuthSubmit(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const authCode = formData.get('authCode');
    
    if (!authCode.trim()) {
        showToast('Please enter an authorization code', 'error');
        return;
    }
    
    showLoading(true);
    
    try {
        const response = await fetch(`${API_BASE}/apply-token`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ authCode: authCode.trim() })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            displayResult('authResult', data);
            showToast('Token applied successfully!', 'success');
            
            // Auto-fill access token in other forms
            if (data.accessToken) {
                autoFillAccessToken(data.accessToken);
            }
        } else {
            displayResult('authResult', data, true);
            showToast(data.error || 'Failed to apply token', 'error');
        }
    } catch (error) {
        displayResult('authResult', { error: error.message }, true);
        showToast('Network error occurred', 'error');
    } finally {
        showLoading(false);
    }
}

// Handle User Info Form Submit
async function handleUserInfoSubmit(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const accessToken = formData.get('accessToken');
    
    if (!accessToken.trim()) {
        showToast('Please enter an access token', 'error');
        return;
    }
    
    showLoading(true);
    
    try {
        const response = await fetch(`${API_BASE}/user-info`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ accessToken: accessToken.trim() })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            displayUserInfo(data);
            showToast('User information retrieved successfully!', 'success');
        } else {
            displayResult('userResult', data, true);
            showToast(data.error || 'Failed to get user info', 'error');
        }
    } catch (error) {
        displayResult('userResult', { error: error.message }, true);
        showToast('Network error occurred', 'error');
    } finally {
        showLoading(false);
    }
}

// Display User Info in a formatted way
function displayUserInfo(data) {
    const resultSection = document.getElementById('userResult');
    
    if (data.userInfo) {
        const userInfo = data.userInfo;
        const userName = userInfo.userName;
        const nameInArabic = userInfo.userNameInArabic;
        
        resultSection.className = 'result-section result-success show';
        resultSection.innerHTML = `
            <div class="result-header">
                <i class="fas fa-check-circle"></i>
                User Information Retrieved
            </div>
            <div class="user-info-card">
                <div class="user-info-header">
                    <div class="user-avatar">
                        ${userInfo.avatar ? `<img src="${userInfo.avatar}" alt="Avatar" style="width: 100%; height: 100%; border-radius: 50%; object-fit: cover;">` : '<i class="fas fa-user"></i>'}
                    </div>
                    <div dir="rtl">
                        <h3>${userName.fullName || 'N/A'}</h3>
                        <p style="color: #718096; margin-top: 4px;">${nameInArabic.fullName || 'N/A'}</p>
                    </div>
                </div>
                <div class="user-details">
                    <div class="user-detail">
                        <span class="user-detail-label">User ID</span>
                        <span class="user-detail-value">${userInfo.userId || 'N/A'}</span>
                    </div>
                    <div class="user-detail">
                        <span class="user-detail-label">Gender</span>
                        <span class="user-detail-value">${userInfo.gender || 'N/A'}</span>
                    </div>
                    <div class="user-detail">
                        <span class="user-detail-label">Birth Date</span>
                        <span class="user-detail-value">${userInfo.birthDate || 'N/A'}</span>
                    </div>
                    <div class="user-detail">
                        <span class="user-detail-label">Nationality</span>
                        <span class="user-detail-value">${userInfo.nationality || 'N/A'}</span>
                    </div>
                    <div class="user-detail">
                        <span class="user-detail-label">Contact Info</span>
                        <span class="user-detail-value">
                            ${userInfo.contactInfos && userInfo.contactInfos.length > 0 
                                ? userInfo.contactInfos.map(contact => `${contact.contactType}: ${contact.contactNo}`).join('<br>')
                                : 'N/A'
                            }
                        </span>
                    </div>
                    <div class="user-detail">
                        <span class="user-detail-label">Login IDs</span>
                        <span class="user-detail-value">
                            ${userInfo.loginIdInfos && userInfo.loginIdInfos.length > 0
                                ? userInfo.loginIdInfos.map(login => `${login.loginIdType}: ${login.maskLoginId}`).join('<br>')
                                : 'N/A'
                            }
                        </span>
                    </div>
                </div>
            </div>
            <div class="result-data" style="margin-top: 1rem;">${JSON.stringify(data, null, 2)}</div>
        `;
    } else {
        displayResult('userResult', data, true);
    }
}

// Handle Cards Form Submit
async function handleCardsSubmit(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const accessToken = formData.get('accessToken');
    
    if (!accessToken.trim()) {
        showToast('Please enter an access token', 'error');
        return;
    }
    
    showLoading(true);
    
    try {
        const response = await fetch(`${API_BASE}/user-cards`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ accessToken: accessToken.trim() })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            displayCardList(data);
            showToast('Card list retrieved successfully!', 'success');
        } else {
            displayResult('cardsResult', data, true);
            showToast(data.error || 'Failed to get card list', 'error');
        }
    } catch (error) {
        displayResult('cardsResult', { error: error.message }, true);
        showToast('Network error occurred', 'error');
    } finally {
        showLoading(false);
    }
}

// Display Card List in a formatted way
function displayCardList(data) {
    const resultSection = document.getElementById('cardsResult');
    
    if (data.cardList) {
        const cards = data.cardList;
        
        resultSection.className = 'result-section result-success show';
        resultSection.innerHTML = `
            <div class="result-header">
                <i class="fas fa-check-circle"></i>
                Card List Retrieved (${cards.length} card${cards.length !== 1 ? 's' : ''})
            </div>
            <div class="card-list">
                ${cards.length > 0 
                    ? cards.map(card => `
                        <div class="payment-card">
                            <div class="card-number">${card.maskedCardNo || 'N/A'}</div>
                            <div class="account-number">Account: ${card.accountNumber || 'N/A'}</div>
                        </div>
                    `).join('')
                    : '<p style="text-align: center; color: #718096; padding: 2rem;">No cards found</p>'
                }
            </div>
            <div class="result-data" style="margin-top: 1rem;">${JSON.stringify(data, null, 2)}</div>
        `;
    } else {
        displayResult('cardsResult', data, true);
    }
}

// Handle Payment Form Submit
async function handlePaymentSubmit(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const paymentData = {
        amount: parseInt(formData.get('amount')),
        requestId: formData.get('requestId').trim(),
        accessToken: formData.get('accessToken').trim(),
        customerId: formData.get('customerId').trim(),
        orderDesc: formData.get('orderDesc').trim(),
        notifyUrl: formData.get('notifyUrl').trim()
    };
    
    // Validation
    if (paymentData.amount <= 0) {
        showToast('Amount must be greater than 0', 'error');
        return;
    }
    
    if (!paymentData.requestId) {
        showToast('Please enter a request ID', 'error');
        return;
    }
    
    if (!paymentData.accessToken) {
        showToast('Please enter an access token', 'error');
        return;
    }
    
    if (!paymentData.customerId) {
        showToast('Please enter a customer ID', 'error');
        return;
    }
    
    if (!paymentData.orderDesc) {
        showToast('Please enter an order description', 'error');
        return;
    }
    
    if (!paymentData.notifyUrl) {
        showToast('Please enter a notification URL', 'error');
        return;
    }
    
    // URL validation
    try {
        new URL(paymentData.notifyUrl);
    } catch {
        showToast('Please enter a valid notification URL', 'error');
        return;
    }
    
    showLoading(true);
    
    try {
        const response = await fetch(`${API_BASE}/pay`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(paymentData)
        });
        
        const data = await response.json();
        
        if (response.ok) {
            displayPaymentResult(data);
            showToast('Payment initiated successfully!', 'success');
            
            // Generate new request ID for next payment
            generateRequestId();
        } else {
            displayResult('paymentResult', data, true);
            showToast(data.error || 'Failed to process payment', 'error');
        }
    } catch (error) {
        displayResult('paymentResult', { error: error.message }, true);
        showToast('Network error occurred', 'error');
    } finally {
        showLoading(false);
    }
}

// Display Payment Result
function displayPaymentResult(data) {
    const resultSection = document.getElementById('paymentResult');
    
    resultSection.className = 'result-section result-success show';
    resultSection.innerHTML = `
        <div class="result-header">
            <i class="fas fa-check-circle"></i>
            Payment Initiated Successfully
        </div>
        <div style="background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 8px; padding: 1.5rem; margin: 1rem 0;">
            <div style="display: grid; gap: 1rem;">
                <div>
                    <strong>Payment ID:</strong> 
                    <span style="font-family: monospace; background: #edf2f7; padding: 2px 6px; border-radius: 4px;">
                        ${data.paymentId || 'N/A'}
                    </span>
                </div>
                <div>
                    <strong>Status:</strong> 
                    <span style="color: #38a169; font-weight: 600;">
                        ${data.result?.resultMessage || 'Success'}
                    </span>
                </div>
                ${data.redirectActionForm?.redirectUrl ? `
                    <div>
                        <strong>Next Step:</strong>
                        <p style="margin: 0.5rem 0; color: #718096;">Redirect user to complete payment:</p>
                        <a href="${data.redirectActionForm.redirectUrl}" target="_blank" 
                           style="display: inline-flex; align-items: center; gap: 8px; padding: 8px 16px; 
                                  background: #f59e0b; color: white; text-decoration: none; border-radius: 6px; font-size: 0.9rem;">
                            <i class="fas fa-external-link-alt"></i>
                            Complete Payment
                        </a>
                    </div>
                ` : ''}
            </div>
        </div>
        <div class="result-data">${JSON.stringify(data, null, 2)}</div>
    `;
}

// Utility function to copy text to clipboard
function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
        showToast('Copied to clipboard!', 'success');
    }).catch(() => {
        showToast('Failed to copy to clipboard', 'error');
    });
}

// Periodically check API health
setInterval(checkAPIHealth, 30000); // Check every 30 seconds
