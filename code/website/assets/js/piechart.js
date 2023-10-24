define(function () {

    class Pie {
        constructor(scores, htmlID) {
            Chart.defaults.color = 'white';
            Chart.defaults.font.family = 'Trispace';
            this.xValues = ["N/A", "Low", "Medium", "High", "Critical"];
            this.yValues = [scores.NA, scores.Low, scores.Medium, scores.High, scores.Critical];
            const barColors = [
                'rgb(100, 100, 100)',
                'rgb(35,206,107)',
                'rgb(54, 162, 235)',
                'rgb(255, 205, 86)',
                'rgb(255, 99, 132)'
            ];

            this.chart = new Chart(htmlID, {
                type: "doughnut",
                data: {
                    labels: this.xValues,
                    datasets: [{
                        backgroundColor: barColors,
                        data: this.yValues
                    }]
                }
            });
        }
        
        resetPie() {
            this.chart.destroy();
        }
    }


    module = {
        "Pie": Pie
    };

    return module;
});