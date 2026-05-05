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
        new Chart(document.getElementById('statusPieChart'), {
            type: 'pie',
            data: {
                labels: status.map(s => s.status),
                datasets: [{
                    data: status.map(s => s.count),
                    backgroundColor: ['#28a745', '#dc3545', '#ffc107', '#17a2b8']
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
