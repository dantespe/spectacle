function get_color() {
    return Math.floor(Math.random()*16777215).toString(16);
}

function color_each(num = 1) {
    return Array(num).fill("#" + get_color())
}

function dim_color(hue = .5, borderColors) {
    var bgColors = [];
    borderColors.forEach((item) => {
        bgColors.push(pSBC(.50, item));
    });
    return bgColors;
}

function create_bar_chart(chart_type, x_input, y_input) {
    var datasets = [];
    y_input.forEach(element => {
        var borderColors = color_each(element.length);    
        var bgColors = dim_color(.5, borderColors=borderColors);
        datasets.push({
            data: element['data'],
            label: element['dataName'],
            borderWidth: 1
        })
    });
    
    const data = {
        labels: x_input[0]['data'],
        datasets: datasets
    };

    if (chart_type == 'stackedBar' || chart_type == 'bar') {
        return {
            type: 'bar',
            data: data,
            options: {
                scales: {
                    y: {
                        beginAtZero: true
                    }
                },
                title: {
                    display: true,
                    text: 'World population per region (in millions)',
                }
            },
        };
    }
    else {
        return {
            type: 'bar',
            data: data,
            options: {
                scales: {
                    y: {
                        beginAtZero: true
                    }
                },
                title: {
                    display: true,
                    text: 'World population per region (in millions)',
                },
                indexAxis: 'y'
            },
        };
    }
}

function create_doughnut_chart(x_input, y_input) {
    var input_data = x_input[0]['data']
    const data = {
        labels: y_input[0]['data'],
        datasets: [{
        label: 'My First Dataset',
        data: input_data,
        backgroundColor: color_each(input_data.length),
        hoverOffset: 4
        }]
    };
    return {
        type: 'doughnut', 
        data: data, 
        options: {
            aspectRatio: 2
        }
    };
}

function create_line_chart(x_input, y_input) {
    const y_data = [];
    y_input.forEach((item) => {
        y_data.push({
            data: item.data,
            label: item.dataName,
            borderColor: "#" + get_color(),
            fill: false
        });
    })
    var data = {
        labels: x_input,
        datasets: y_data
    };

    return {
        type: 'line', 
        data: data, 
        options: {}
    };
}

function create_scatter_chart(x_input, y_input) {
    const data = {
        datasets: [{
        label: 'Scatter Dataset',
        data: [{
            x: -10,
            y: 0
        }, {
            x: 0,
            y: 10
        }, {
            x: 10,
            y: 5
        }, {
            x: 0.5,
            y: 5.5
        }],
        backgroundColor: 'rgb(255, 99, 132)'
        }],
    };
    return {
        type: 'scatter', 
        data: data, 
        options: {}
    };
}

function create_bubble_chart(x_input, y_input) {
    const data = {
        datasets: [{
        label: 'First Dataset',
        data: [{
            x: 20,
            y: 30,
            r: 15
        }, {
            x: 40,
            y: 10,
            r: 10
        }],
        backgroundColor: 'rgb(255, 99, 132)'
        }]
    };
    return {
        type: 'bubble', 
        data: data, 
        options: {}
    };
}

function create_pie_chart(x_input, y_input) {
    const data = {
        labels: x_input[0]['data'],
        datasets: [{
        label: 'My First Dataset',
        data: y_input[0]['data'],
        backgroundColor: [
            'rgb(255, 99, 132)',
            'rgb(54, 162, 235)',
            'rgb(255, 205, 86)'
        ],
        hoverOffset: 4
        }]
    };
    return {
        type: 'pie', 
        data: data, 
        options: {
            aspectRatio: 2
        }
    };
}

function create_radar_chart(x_input, y_input) {
    const data = {
        labels: [
            'Eating',
            'Drinking',
            'Sleeping',
            'Designing',
            'Coding',
            'Cycling',
            'Running'
        ],
        datasets: [{
            label: 'My First Dataset',
            data: [65, 59, 90, 81, 56, 55, 40],
            fill: true,
            backgroundColor: 'rgba(255, 99, 132, 0.2)',
            borderColor: 'rgb(255, 99, 132)',
            pointBackgroundColor: 'rgb(255, 99, 132)',
            pointBorderColor: '#fff',
            pointHoverBackgroundColor: '#fff',
            pointHoverBorderColor: 'rgb(255, 99, 132)'
        }, {
            label: 'My Second Dataset',
            data: [28, 48, 40, 19, 96, 27, 100],
            fill: true,
            backgroundColor: 'rgba(54, 162, 235, 0.2)',
            borderColor: 'rgb(54, 162, 235)',
            pointBackgroundColor: 'rgb(54, 162, 235)',
            pointBorderColor: '#fff',
            pointHoverBackgroundColor: '#fff',
            pointHoverBorderColor: 'rgb(54, 162, 235)'
        }]
    };
    return {
        type: "radar", 
        data: data, 
        options: {
            elements: {
                line: {
                    borderWidth: 3
                }
            },
            aspectRatio: 2
        }
    };
}

function create_polarArea_chart(x_input, y_input) {
    const data = {
        labels: y_input[0]['data'],
        datasets: [{
            label: 'My First Dataset',
            data: x_input[0]['data'],
            backgroundColor: [
                'rgb(255, 99, 132)',
                'rgb(75, 192, 192)',
                'rgb(255, 205, 86)',
                'rgb(201, 203, 207)',
                'rgb(54, 162, 235)'
            ]
        }]
    };
    return {
        type: 'polarArea', 
        data: data, 
        options: {
            aspectRatio: 2
        }
    };
}

const pSBC=(p,c0,c1,l)=>{
    let r,g,b,P,f,t,h,i=parseInt,m=Math.round,a=typeof(c1)=="string";
    if(typeof(p)!="number"||p<-1||p>1||typeof(c0)!="string"||(c0[0]!='r'&&c0[0]!='#')||(c1&&!a))return null;
    if(!this.pSBCr)this.pSBCr=(d)=>{
        let n=d.length,x={};
        if(n>9){
            [r,g,b,a]=d=d.split(","),n=d.length;
            if(n<3||n>4)return null;
            x.r=i(r[3]=="a"?r.slice(5):r.slice(4)),x.g=i(g),x.b=i(b),x.a=a?parseFloat(a):-1
        }else{
            if(n==8||n==6||n<4)return null;
            if(n<6)d="#"+d[1]+d[1]+d[2]+d[2]+d[3]+d[3]+(n>4?d[4]+d[4]:"");
            d=i(d.slice(1),16);
            if(n==9||n==5)x.r=d>>24&255,x.g=d>>16&255,x.b=d>>8&255,x.a=m((d&255)/0.255)/1000;
            else x.r=d>>16,x.g=d>>8&255,x.b=d&255,x.a=-1
        }return x};
    h=c0.length>9,h=a?c1.length>9?true:c1=="c"?!h:false:h,f=this.pSBCr(c0),P=p<0,t=c1&&c1!="c"?this.pSBCr(c1):P?{r:0,g:0,b:0,a:-1}:{r:255,g:255,b:255,a:-1},p=P?p*-1:p,P=1-p;
    if(!f||!t)return null;
    if(l)r=m(P*f.r+p*t.r),g=m(P*f.g+p*t.g),b=m(P*f.b+p*t.b);
    else r=m((P*f.r**2+p*t.r**2)**0.5),g=m((P*f.g**2+p*t.g**2)**0.5),b=m((P*f.b**2+p*t.b**2)**0.5);
    a=f.a,t=t.a,f=a>=0||t>=0,a=f?a<0?t:t<0?a:a*P+t*p:0;
    if(h)return"rgb"+(f?"a(":"(")+r+","+g+","+b+(f?","+m(a*1000)/1000:"")+")";
    else return"#"+(4294967296+r*16777216+g*65536+b*256+(f?m(a*255):0)).toString(16).slice(1,f?undefined:-2)
}