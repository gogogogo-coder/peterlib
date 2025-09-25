package peterlib

import (
	"fmt"
	"os"
)

type Command struct {
	Cmd   string
	Value interface{}
}

var log []Command
const canvasWidth = 400
const canvasHeight = 400

// -------------------------------------------------
// Fonctions utilisables dans le code Go
// -------------------------------------------------
func Down()            { log = append(log, Command{"Down", nil}) }
func Up()              { log = append(log, Command{"Up", nil}) }
func Right()           { log = append(log, Command{"Right", nil}) }
func Left()            { log = append(log, Command{"Left", nil}) }
func Pivote(angle int) { log = append(log, Command{"Pivote", angle}) }
func Forward(n int)    { log = append(log, Command{"Forward", n}) }
func Color(c string)   { log = append(log, Command{"Color", c}) }
func Say(msg string)   { log = append(log, Command{"Say", msg}) }

// -------------------------------------------------
// Génération automatique du HTML à la fin
// -------------------------------------------------
func init() {
	// La fonction generateHTML() sera appelée automatiquement à la fin du programme
	done := make(chan bool)
	go func() {
		<-done
		if len(log) > 0 {
			generateHTML()
		}
	}()
	// Hook pour fermer le channel à la fin du main()
	defer func() { done <- true }()
}

// -------------------------------------------------
// Génération du fichier HTML dynamique
// -------------------------------------------------
func generateHTML() {
	f, _ := os.Create("peter.html")
	defer f.Close()

	fmt.Fprintln(f, `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Peter le traceur</title>
<style>
canvas { border:1px solid black; background:#f9f9f9; }
</style>
</head>
<body>
<canvas id="c" width="400" height="400"></canvas>
<script>
const cmds = [`)

	for i, c := range log {
		if c.Value != nil {
			fmt.Fprintf(f, `{cmd:"%s", value:%q}`, c.Cmd, c.Value)
		} else {
			fmt.Fprintf(f, `{cmd:"%s"}`, c.Cmd)
		}
		if i < len(log)-1 {
			fmt.Fprintln(f, ",")
		}
	}

	fmt.Fprintf(f, `];

const canvas = document.getElementById("c");
const ctx = canvas.getContext("2d");
let x = %d, y = %d, angle = 0, penDown = false;
ctx.lineWidth = 2;
ctx.strokeStyle = "black";

function forward(n){
  let rad = angle * Math.PI/180;
  let x2 = x + n * Math.cos(rad);
  let y2 = y + n * Math.sin(rad);
  if(penDown){
    ctx.beginPath();
    ctx.moveTo(x, y);
    ctx.lineTo(x2, y2);
    ctx.stroke();
  }
  x = x2; y = y2;
  ctx.beginPath();
  ctx.arc(x, y, 3, 0, Math.PI*2);
  ctx.fillStyle = "red";
  ctx.fill();
}

let i=0;
function step(){
  if(i >= cmds.length) return;
  let c = cmds[i++];
  switch(c.cmd){
    case "Down": penDown=true; break;
    case "Up": penDown=false; break;
    case "Right": angle += 90; break;
    case "Left": angle -= 90; break;
    case "Pivote": angle += c.value; break;
    case "Forward": forward(c.value); break;
    case "Color": ctx.strokeStyle=c.value; break;
    case "Say": alert("Peter dit: " + c.value); break;
  }
}

setInterval(step,500);
</script>
</body>
</html>`, canvasWidth/2, canvasHeight/2)
}
