package gui

import (
    "fmt"
    "image/color"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/widget"

    "Estacionamiento/internal/domain/models"
)

type GUI struct {
    parkingSpots    [20]*canvas.Image
    parkingSpaces   [20]*fyne.Container
    window          fyne.Window
    statsLabels     map[string]*widget.Label
    mainContainer   *fyne.Container
    gateRect        *canvas.Rectangle
}

func createParkingSpace() *canvas.Image {
    img := canvas.NewImageFromFile("")
    img.SetMinSize(fyne.NewSize(100, 100))
    img.FillMode = canvas.ImageFillContain
    img.Resize(fyne.NewSize(100, 100))
    return img
}

func createStatsLabel(text string) *widget.Label {
    label := widget.NewLabel(text)
    label.Alignment = fyne.TextAlignCenter
    return label
}

func NewGUI() *GUI {
    myApp := app.New()
    window := myApp.NewWindow("Simulador de Estacionamiento")
    
    gui := &GUI{
        window:      window,
        statsLabels: make(map[string]*widget.Label),
    }

    gui.statsLabels["generated"] = createStatsLabel("Carros generados: 0/100")
    gui.statsLabels["processed"] = createStatsLabel("Carros procesados: 0")
    gui.statsLabels["parked"] = createStatsLabel("Carros estacionados: 0/20")
    gui.statsLabels["waitingEnter"] = createStatsLabel("Cola para entrar: 0")
    gui.statsLabels["waitingExit"] = createStatsLabel("Cola para salir: 0")

    statsContainer := container.NewVBox(
        gui.statsLabels["generated"],
        gui.statsLabels["processed"],
        gui.statsLabels["parked"],
        gui.statsLabels["waitingEnter"],
        gui.statsLabels["waitingExit"],
    )

    gui.gateRect = canvas.NewRectangle(color.RGBA{R: 0, G: 255, B: 0, A: 255})
    gui.gateRect.Resize(fyne.NewSize(80, 20))
    gateContainer := container.NewHBox(
        layout.NewSpacer(),
        gui.gateRect,
        layout.NewSpacer(),
    )

    leftBorder := canvas.NewRectangle(color.RGBA{R: 180, G: 180, B: 180, A: 255})
    leftBorder.Resize(fyne.NewSize(10, 500))
    rightBorder := canvas.NewRectangle(color.RGBA{R: 180, G: 180, B: 180, A: 255})
    rightBorder.Resize(fyne.NewSize(10, 500))

    parkingGrid := container.New(layout.NewGridLayout(4))
    
    for i := 0; i < 20; i++ {
        divider := canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 200, A: 255})
        divider.Resize(fyne.NewSize(2, 100))
        
        emptySpace := createParkingSpace()
        gui.parkingSpots[i] = emptySpace
        
        spaceContainer := container.NewHBox(
            divider,
            emptySpace,
        )
        gui.parkingSpaces[i] = spaceContainer
        parkingGrid.Add(spaceContainer)
    }

    parkingContainer := container.NewHBox(
        leftBorder,
        parkingGrid,
        rightBorder,
    )

    gui.mainContainer = container.NewVBox(
        statsContainer,
        gateContainer,
        parkingContainer,
    )

    window.SetContent(container.NewPadded(gui.mainContainer))
    window.Resize(fyne.NewSize(900, 700))
    
    return gui
}

func (gui *GUI) Update(parkingSpot int, isOccupied bool) {
    if isOccupied {
        newImg := canvas.NewImageFromFile("assets/carroRojo.png")
        newImg.SetMinSize(fyne.NewSize(100, 100))
        newImg.FillMode = canvas.ImageFillContain
        newImg.Resize(fyne.NewSize(100, 100))
        gui.parkingSpots[parkingSpot] = newImg
        gui.parkingSpaces[parkingSpot].Objects[1] = newImg
    } else {
        emptySpace := createParkingSpace()
        gui.parkingSpots[parkingSpot] = emptySpace
        gui.parkingSpaces[parkingSpot].Objects[1] = emptySpace
    }
    gui.parkingSpaces[parkingSpot].Refresh()
}

func (gui *GUI) UpdateGate(isOccupied bool) {
    if isOccupied {
        gui.gateRect.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
    } else {
        gui.gateRect.FillColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
    }
    gui.gateRect.Refresh()
}

func (gui *GUI) UpdateStats(stats models.ParkingStats) {
    gui.statsLabels["generated"].SetText(fmt.Sprintf("Carros generados: %d/100", stats.TotalGenerated))
    gui.statsLabels["processed"].SetText(fmt.Sprintf("Carros procesados: %d", stats.TotalProcessed))
    gui.statsLabels["parked"].SetText(fmt.Sprintf("Carros estacionados: %d/20", stats.CurrentlyParked))
    gui.statsLabels["waitingEnter"].SetText(fmt.Sprintf("Cola para entrar: %d", stats.WaitingToEnter))
    gui.statsLabels["waitingExit"].SetText(fmt.Sprintf("Cola para salir: %d", stats.WaitingToExit))
}

func (gui *GUI) GetWindow() fyne.Window {
    return gui.window
}