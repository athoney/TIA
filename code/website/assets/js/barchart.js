define(function () {
    
    class Bar {
        constructor(cweId, counts, htmlID, cwename) {
            Chart.defaults.color = 'white';
            Chart.defaults.font.family = 'Trispace';
            this.labels = cweId;
            this.data = counts;
            const barColors = [
                'rgb(255, 99, 132)',
                'rgb(54, 162, 235)',
                'rgb(255, 206, 86)',
                'rgb(75, 192, 192)',
                'rgb(153, 102, 255)',
                'rgb(255, 159, 64)',
                'rgb(200, 50, 132)',
                'rgb(4, 62, 235)',
                'rgb(155, 206, 86)',
                'rgb(10, 192, 255)',
                'rgb(200, 102, 255)',
                'rgb(100, 200, 100)'
            ];
            this.cwename = cwename;
            this.chart = new Chart(htmlID, {
                type: 'bar',
                data: {
                    labels: this.labels,
                    datasets: [{
                        axis: 'y',
                        label: 'count',
                        data: this.data,
                        backgroundColor: barColors,
                        borderColor: barColors,
                        borderWidth: 1
                    }]
                },
                options: {
                    indexAxis: 'y',
                    xAxes: [{
                        ticks: {
                            min: 0 // Edit the value according to what you need
                        }
                    }],
                    plugins: {
                        tooltip: {
                            callbacks: {
                                title: function(context){
                                    return cwename[context[0].dataIndex];
                                }
                            }
                        }
                    }
                }
            });
            window.myBarChart = this.chart;
        }

        resetBar() {
            this.chart.destroy();
        }
    }

    module = {
        "Bar": Bar
    };

    return module;
});