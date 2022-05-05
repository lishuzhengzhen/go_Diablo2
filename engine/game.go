package engine

import (
	"embed"
	"game/controller"
	"game/fonts"
	"game/interfaces"
	"game/layout"
	"game/mapCreator/mapManage"
	"game/music"
	"game/role"
	"game/status"
	"game/storage"
	"game/tools"
	"runtime"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//配置信息
const (
	SCREENWIDTH  int = 490
	SCREENHEIGHT int = 300
	WEOFFSETX    int = 127
	WEOFFSETY    int = 14
)

type Game struct {
	count, countForMap int
	player             [2]*role.Player           //玩家
	mapManage          interfaces.MapInterface   //地图等管理
	ui                 *layout.UI                //UI
	music              interfaces.MusicInterface //音乐
	status             *status.StatusManage      //状态管理器
	font_style         *fonts.FontBase           //字体
}

var (
	counts      int = 0
	countsFor20 int = 0
	countsFor8  int = 0
	countsFor17 int = 0
	frameNums   int = 4
	frameSpeed  int = 5
	mouseX      int
	mouseY      int
)

//go:embed resource
var asset embed.FS

//GameEngine
func NewGame() *Game {
	// w := ws.NewNet()
	// w.Start()
	//statueManage
	sta := status.NewStatusManage()
	bag := storage.New()
	//场景
	m := mapManage.NewE1(&asset, sta, bag)
	//Player  设置初始状态和坐标
	r := role.NewPlayer(5280, 1880, tools.IDLE, 0, 0, 0, &asset, m, sta)
	r1 := role.NewPlayer(5280, 1880, tools.IDLE, 0, 0, 0, &asset, m, sta)

	//字体
	f := fonts.NewFont(&asset)
	//UI
	u := layout.NewUI(&asset, sta, f, m, bag)
	bag.UI = u
	//BGM
	bgm := music.NewMusicBGM(&asset)

	gameEngine := &Game{
		count:       0,
		countForMap: 0,
		player:      [2]*role.Player{r, r1},
		ui:          u,
		music:       bgm,
		status:      sta,
		mapManage:   m,
		font_style:  f,
	}
	//启动游戏
	gameEngine.StartEngine()
	return gameEngine
}

//引擎启动
func (g *Game) StartEngine() {
	//隐藏鼠标系统的ICON
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	w := sync.WaitGroup{}
	w.Add(1)
	//UI Init
	go func() {
		//加载字体
		g.font_style.LoadFont("resource/font/pf_normal.ttf")
		g.ui.LoadGameLoginImages()
		runtime.GC()
		w.Done()
	}()
	w.Wait()
	go func() {
		runtime.GC()
	}()
}

func (g *Game) Update() error {
	//判断是否是点击屏幕
	if !controller.IsTouch() {
		mouseX, mouseY = ebiten.CursorPosition()
	} else {
		mouseX, mouseY = controller.GetTouchDefaultXY()
	}

	//切换场景逻辑
	if !g.status.ChangeScenceFlg {
		switch g.status.CurrentGameScence {
		case tools.GAMESCENESTART:
			//进入游戏场景逻辑
			g.changeScenceGameUpdate()
		case tools.GAMESCENEOPENDOOR:
			//游戏加载逻辑
			g.ChangeScenceOpenDoorUpdate()
		case tools.GAMESCENESELECTROLE:
			//进入游戏选择界面逻辑
			g.ChangeScenceSelectUpdate()
		default:
			//进入游戏登录界面逻辑
			g.ChangeScenceLoginUpdate()
		}
	}
	//全屏显示控制
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		i := ebiten.IsFullscreen()
		ebiten.SetFullscreen(!i)
	}
	//Debug 信息显示控制
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.status.DisPlayDebugInfo = !g.status.DisPlayDebugInfo
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//判断是否切换场景
	if !g.status.ChangeScenceFlg {
		switch g.status.CurrentGameScence {
		case tools.GAMESCENESTART:
			g.ChangeScenceGameDraw(screen)
		case tools.GAMESCENESELECTROLE:
			g.ChangeScenceSelectDraw(screen)
		case tools.GAMESCENEOPENDOOR:
			g.ChangeScenceOpenDoorDraw(screen)
		default:
			g.ChangeScenceLoginDraw(screen)
		}
	}
	//绘制鼠标ICON
	g.ui.DrawMouseIcon(screen, mouseX, mouseY)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return tools.LAYOUTX, tools.LAYOUTY
}
