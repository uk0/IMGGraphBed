var Canvas = document.getElementById("myCanvas");
var can = Canvas.getContext("2d");
Canvas.width = document.documentElement.clientWidth;
Canvas.height = document.documentElement.clientHeight;
var dianArr = [];

function randNum(min, max) {
    return parseInt(Math.random() * (max - min + 1) + min);
}
function Dian() {
    this.x = randNum(Canvas.width / 2 - 200, (Canvas.width / 2 + 200));
    this.y = randNum(Canvas.height / 2 - 200, (Canvas.height / 2 + 200));
    this.w = randNum(1, 5);
    this.h = randNum(1, 5);
    this.speedX = 0;
    this.speedY = 0;
    this.bluX = randNum(0, 1); //Y正负
    this.bluY = randNum(0, 1); //X正负
    this.blu = randNum(1, 2); // 将所有点分为两类；
    this.blu1 = randNum(1, 2); //正方向点XY正负；
}
Dian.prototype.move = function() {
    if (this.blu == 1) {
        if (this.blu1 == 1) {
            if (this.bluX == 1) {
                this.speedX += Math.random() / 250;
                this.x += this.speedX;
            } else {
                this.speedX += Math.random() / 250;
                this.x -= this.speedX;
            }
            if (this.bluY == 1) {
                this.speedY += Math.random() / 50;
                this.y -= this.speedY;
            } else {
                this.speedY += Math.random() / 50;
                this.y += this.speedY;
            }
        } else {
            if (this.bluX == 1) {
                this.speedX += Math.random() / 25;
                this.x += this.speedX;
            } else {
                this.speedX += Math.random() / 25;
                this.x -= this.speedX;
            }
            if (this.bluY == 1) {
                this.speedY += Math.random() / 400;
                this.y -= this.speedY;
            } else {
                this.speedY += Math.random() / 400;
                this.y += this.speedY;
            }
        }
    } else {
        this.speedY += Math.random() / 50;
        this.speedX += Math.random() / 25;
        if (this.bluX == 1) {
            this.x += this.speedX;
        } else {
            this.x -= this.speedX;
        }
        if (this.bluY == 1) {
            this.y -= this.speedY;
        } else {
            this.y += this.speedY;
        }
    }
}
Dian.prototype.drawDian = function() {
    can.fillStyle = "white";
    can.fillRect(this.x, this.y, this.w, this.h);
}
function drawBGC() {
    can.clearRect(0, 0, Canvas.width, Canvas.height);
    can.beginPath();
    can.fillStyle = "black";
    can.fillRect(0, 0, Canvas.width, Canvas.height);
    for (var i = 0; i < 10; i++) {
        var dian = new Dian();
        dianArr.push(dian);
    }
    for (var j = 0; j < dianArr.length; j++) {
        var aa = dianArr[j];
        aa.move();
        aa.drawDian();
        if (aa.x > Canvas.width || aa.x < 0 || aa.y > Canvas.height || aa.y < 0) {
            dianArr.splice(i, 1);
        }
    }
    can.closePath();
    window.requestAnimationFrame(drawBGC);
}
drawBGC();