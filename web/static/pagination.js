(function() {
    const rowsPerPage = 20;
    const table = document.querySelector('table');
    if (!table) return;

    const rows = table.querySelectorAll('tr');
    const dataRows = Array.from(rows).filter(row => row.querySelector('td'));
    if (dataRows.length <= rowsPerPage) return;

    let currentPage = 1;
    const totalPages = Math.ceil(dataRows.length / rowsPerPage);

    dataRows.forEach(row => (row.style.display = 'none'));

    function showPage(page) {
        const start = (page - 1) * rowsPerPage;
        const end = start + rowsPerPage;
        dataRows.forEach((row, index) => {
            row.style.display = (index >= start && index < end) ? '' : 'none';
        });
    }

    const paginationContainer = document.createElement('div');
    paginationContainer.style.cssText = 'text-align:center; margin-top:1rem;';

    function renderButtons() {
        paginationContainer.innerHTML = '';

        const prevBtn = document.createElement('button');
        prevBtn.textContent = 'Previous';
        prevBtn.disabled = currentPage === 1;
        prevBtn.style.cssText = 'margin:0 0.3rem; padding:0.3rem 0.7rem;';
        prevBtn.addEventListener('click', () => {
            if (currentPage > 1) {
                currentPage--;
                showPage(currentPage);
                renderButtons();
            }
        });
        paginationContainer.appendChild(prevBtn);

        const maxVisible = 5;
        let startPage = Math.max(1, currentPage - 2);
        let endPage = Math.min(totalPages, startPage + maxVisible - 1);
        if (endPage - startPage + 1 < maxVisible) {
            startPage = Math.max(1, endPage - maxVisible + 1);
        }

        if (startPage > 1) {
            const firstBtn = document.createElement('button');
            firstBtn.textContent = '1';
            firstBtn.style.cssText = 'margin:0 0.3rem; padding:0.3rem 0.7rem;';
            firstBtn.addEventListener('click', () => { currentPage = 1; showPage(currentPage); renderButtons(); });
            paginationContainer.appendChild(firstBtn);
            if (startPage > 2) {
                const dots = document.createElement('span');
                dots.textContent = '...';
                dots.style.margin = '0 0.3rem';
                paginationContainer.appendChild(dots);
            }
        }

        for (let p = startPage; p <= endPage; p++) {
            const pageBtn = document.createElement('button');
            pageBtn.textContent = p;
            pageBtn.style.cssText = 'margin:0 0.3rem; padding:0.3rem 0.7rem;';
            if (p === currentPage) {
                pageBtn.style.fontWeight = 'bold';
                pageBtn.style.backgroundColor = '#1a73e8';
                pageBtn.style.color = 'white';
            }
            pageBtn.addEventListener('click', () => {
                currentPage = p;
                showPage(currentPage);
                renderButtons();
            });
            paginationContainer.appendChild(pageBtn);
        }

        if (endPage < totalPages) {
            if (endPage < totalPages - 1) {
                const dots = document.createElement('span');
                dots.textContent = '...';
                dots.style.margin = '0 0.3rem';
                paginationContainer.appendChild(dots);
            }
            const lastBtn = document.createElement('button');
            lastBtn.textContent = totalPages;
            lastBtn.style.cssText = 'margin:0 0.3rem; padding:0.3rem 0.7rem;';
            lastBtn.addEventListener('click', () => { currentPage = totalPages; showPage(currentPage); renderButtons(); });
            paginationContainer.appendChild(lastBtn);
        }

        const nextBtn = document.createElement('button');
        nextBtn.textContent = 'Next';
        nextBtn.disabled = currentPage === totalPages;
        nextBtn.style.cssText = 'margin:0 0.3rem; padding:0.3rem 0.7rem;';
        nextBtn.addEventListener('click', () => {
            if (currentPage < totalPages) {
                currentPage++;
                showPage(currentPage);
                renderButtons();
            }
        });
        paginationContainer.appendChild(nextBtn);
    }

    showPage(1);
    renderButtons();
    table.parentNode.insertBefore(paginationContainer, table.nextSibling);
})();
