async function fetchJSON(url) {
    const resp = await fetch(url, { credentials: 'same-origin' });
    if (!resp.ok) throw new Error(`HTTP ${resp.status}`);
    return resp.json();
}

(async () => {
    try {
        const perDay = await fetchJSON('/analytics/artifacts-per-day');
        new Chart(document.getElementById('artifactsPerDayChart'), {
            type: 'bar',
            data: {
                labels: perDay.map(d => d.date),
                datasets: [{
                    label: 'Artifacts per day',
                    data: perDay.map(d => d.count),
                    backgroundColor: 'rgba(54, 162, 235, 0.5)'
                }]
            },
            options: {
                responsive: true,
                plugins: {
                    title: {
                        display: true,
                        text: 'Artifacts per day'
                    }
                }
            }
        });
    } catch (e) {
        console.error('Error: could not load analytics', e);
    }

    try {
        const status = await fetchJSON('/analytics/status-distribution');
        const order = ['Passed', 'Failed', 'Skipped', 'Error', 'Missing'];
        status.sort((a, b) => {
            const indexA = order.indexOf(a.status);
            const indexB = order.indexOf(b.status);

            const ia = indexA === -1 ? order.length : indexA;
            const ib = indexB === -1 ? order.length : indexB;
            return ia - ib;
        });

        const colorMap = {
            'Passed': '#28a745',
            'Failed': '#dc3545',
            'Skipped': '#ffc107',
            'Error': '#5e5e5e',
            'Missing': '#970290'
        };
        const backgroundColors = status.map(s => colorMap[s.status] || '#6c757d');

        new Chart(document.getElementById('statusPieChart'), {
            type: 'pie',
            data: {
                labels: status.map(s => s.status),
                datasets: [{
                    data: status.map(s => s.count),
                    backgroundColor: backgroundColors
                }]
            },
            options: {
                responsive: true,
                plugins: {
                    title: {
                        display: true,
                        text: 'Status distribution'
                    }
                }
            }
        });
    } catch (e) {
        console.error('Error: could not load analytics', e);
    }
})();
