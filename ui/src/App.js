import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css';
import ReactFlow from 'react-flow-renderer';

import citypos from './citypos.json'
import React from 'react';
import { ButtonGroup, Dropdown, Button, DropdownButton } from 'react-bootstrap';


const textColorMap = {
  Black: "white",
  Blue: "white",
  Yellow: "black",
  Red: "black",
}

const xMap = {
  Black: 0,
  Blue: 4,
  Yellow: 0,
  Red: 4,
}
const yMap = {
  Black: 2,
  Blue: 2,
  Yellow: -1,
  Red: -1,
}

const pMap = {
  P0: "White",
  P1: "Green",
  P2: "Pink",
  P3: "Orange",
}

const pPos = {
  P0: 0,
  P1: 1,
  P2: 2,
  P3: 3,
}


const BACKEND = "http://localhost:8080"

class App extends React.Component {
  constructor(props) {
    super(props)

    let elements = [];
    this.state = {
      elements,
      game: null
    }
    this.getCity = this.getCity.bind(this)
    this.onLoad = this.onLoad.bind(this)
    this.click = this.click.bind(this)
    this.updateState = this.updateState.bind(this)
  }

  click(e) {
    const [action, target] = e.split("-")
    const player = this.state.game.State.Players[this.state.game.State.CurrentPlayerN]

    console.log(action, target)
    return fetch(`${BACKEND}/action/${action}?player=${player.Name}&target=${target}`, {
      method: "POST"
    })
      .then(res => res.json())
      .then(res => this.updateState(res))
  }

  getCity(name, elements = this.state.elements) {
    return elements.find(c => c.id === name)
  }

  updateCity(name, update) {
    return this.setState({
      elements: this.state.elements.map((city) => {
        if (city.id === name) {
          city = { ...city, ...update }
        }
        return city
      })
    })
  }

  updateState(input) {
    const elements = []
    const state = input.State
    state.Cities.forEach((city) => {
      const cityX = citypos[city.Name].x * 150
      const cityY = citypos[city.Name].y * 100
      elements.push({
        id: city.Name,
        data: { label: city.Name, Virus: city.VirusType },
        style: {
          width: "80px",
          background: city.VirusType,
          color: textColorMap[city.VirusType],
        },
        position: { x: cityX, y: cityY },
      })
      city.Links.forEach((link) => {
        if (!elements.some(c => c.id === `e${link}-${city.Name}`)) {
          elements.push({
            id: `e${city.Name}-${link}`,
            source: city.Name,
            target: link,
            type: 'straight',
            animated: true,
          })
        }
      })
      if (city.Buildings.ResearchBuilding) {
        elements.push({
          id: `r${city.Name}`,
          data: { label: "" },
          style: {
            width: "1px",
            height: "1px",
            background: "white",
          },
          position: {
            x: cityX - 20,
            y: cityY + 5,
          },
        })
      }

      Object.keys(city.VirusCounts).forEach((vName) => {
        for (let i = 0; i < city.VirusCounts[vName]; i++) {
          elements.push({
            id: `v${city.Name}-${vName}-${i}`,
            data: { label: "" },
            style: {
              width: "20px",
              height: "20px",
              background: vName,
            },
            position: {
              x: cityX - 15 + 15 * xMap[vName] + 15 * i,
              y: cityY + 15 * yMap[vName],
            },
          })
        }
      })
    })

    state.Players.forEach((player) => {
      const pName = player.Name
      const cityName = player.Location
      const cityX = citypos[cityName].x * 150
      const cityY = citypos[cityName].y * 100
      elements.push({
        id: `p${pName}`,
        data: { label: "" },
        style: {
          width: "1px",
          height: "1px",
          background: pMap[pName],
        },
        position: {
          x: cityX + 80,
          y: cityY - 10 + 10 * pPos[pName],
        },
      })

    })


    this.setState({ elements, game: input })
  }

  onLoad(reactFlowInstance) {
    return fetch(`${BACKEND}/state`)
      .then(res => res.json())
      .then(res => this.updateState(res))
      .then(() => setTimeout(() => {
        reactFlowInstance.fitView()
      }, 1))
      .catch(e => console.log("error", e))
  }

  moveButton(playerActions) {
    if (!playerActions.move) return null
    return (
      <DropdownButton
        size="lg"
        drop="down"
        title="move"
        as={ButtonGroup}
        onSelect={this.click}
      >
        {playerActions.move.map((city) => {
          return <Dropdown.Item key={city} eventKey={`move-${city}`}>{city}</Dropdown.Item>
        })}
      </DropdownButton>
    )
  }

  cureButton(playerActions) {
    if (!playerActions.cure) return null

    return (
      <DropdownButton
        size="lg"
        drop="down"
        title="cure"
        disabled={!playerActions.cure}
        as={ButtonGroup}
        onSelect={this.click}
      >
        {playerActions.cure && playerActions.cure.map((city) => {
          return <Dropdown.Item key={city} eventKey={`cure-${city}`}>{city}</Dropdown.Item>
        })}
      </DropdownButton>
    )
  }

  infectButton(playerActions) {
    if (!playerActions.infect) return null

    return (
      <Button
        onClick={() => this.click("infect-self")}
        variant="success"
        size="lg"
        title="infect"
        disabled={!playerActions.infect}
        as={ButtonGroup}
      >
        infect
      </Button>
    )
  }

  epidemicButton(playerActions) {
    if (!playerActions.epidemic) return null
    return (
      <Button
        onClick={() => this.click("epidemic-self")}
        variant="success"
        size="lg"
        title="epidemic"
        disabled={!playerActions.epidemic}
        as={ButtonGroup}
      >
        epidemic
      </Button>
    )
  }

  drawButton(playerActions) {
    if (!playerActions.draw) return null

    return (
      <Button
        onClick={() => this.click("draw-self")}
        size="lg"
        title="draw"
        disabled={!playerActions.draw}
        as={ButtonGroup}
      >
        draw
      </Button>
    )
  }

  outbreakButton(playerActions) {
    if (!playerActions.outbreak) return null
    return (
      <Button
        onClick={() => this.click(`outbreak-${playerActions.outbreak}`)}
        variant="warning"
        size="lg"
        title="outbreak"
        disabled={!playerActions.outbreak}
        as={ButtonGroup}
      >
        outbreak
      </Button>
    )
  }

  researchButton(playerActions) {
    if (!playerActions.research) return null

    return (
      <DropdownButton
        size="lg"
        drop="down"
        title="research"
        disabled={!playerActions.research}
        as={ButtonGroup}
        onSelect={this.click}
      >
        {playerActions.research && playerActions.research.map((city) => {
          return <Dropdown.Item key={city} eventKey={`research-${city}`}>{city}</Dropdown.Item>
        })}
      </DropdownButton>
    )
  }

  getButton(playerActions) {
    if (!playerActions.get) return null

    return (
      <DropdownButton
        size="lg"
        drop="down"
        title="get"
        disabled={!playerActions.get}
        as={ButtonGroup}
        onSelect={this.click}
      >
        {playerActions.get && playerActions.get.map((city) => {
          return <Dropdown.Item key={city} eventKey={`get-${city}`}>{city}</Dropdown.Item>
        })}
      </DropdownButton>
    )
  }
  hand(gameState, player, playerActions) {
    return player.Hand.Cards.map((card) => {
      return (
          <DropdownButton
            variant={card.VirusType.toLowerCase()}
            size="lg"
            drop="up"
            key={card.Name}
            as={ButtonGroup}
            title={card.Name}
            onSelect={this.click}
          >
            {playerActions.flyTo && playerActions.flyTo.includes(card.Name) ? <Dropdown.Item key={`flyTo-${card.Name}`} eventKey={`flyTo-${card.Name}`}>Fly</Dropdown.Item> : null}
            {playerActions.build && playerActions.build.includes(`${card.Name}:ResearchBuilding`) ? <Dropdown.Item key={`build-${card.Name}`} eventKey={`build-${card.Name}:ResearchBuilding`}>Build</Dropdown.Item> : null}
            {playerActions.flyAnywhere && playerActions.flyAnywhere.includes(card.Name) ? <Dropdown.Item key={`flyAnywhere-${card.Name}`} eventKey={`flyAnywhere-${card.Name}`}>Build</Dropdown.Item> : null}
            {playerActions.discard && playerActions.discard.includes(card.Name) ? <Dropdown.Item key={`discard-${card.Name}`} eventKey={`discard-${card.Name}`}>Discard</Dropdown.Item> : null}
            {playerActions.give && playerActions.give.includes(card.Name) ? <Dropdown.Item key={`give-${card.Name}`} eventKey={`give-${card.Name}`}>Give</Dropdown.Item> : null}
          </DropdownButton>
      )
    })
  }
  actions() {
    if (!this.state.game) return null
    const gameState = this.state.game.State
    const player = gameState.Players[gameState.CurrentPlayerN]
    const playerActions = this.state.game.Actions[player.Name]

    return (
      <div>
        <h1>Player {player.Name}: Actions {gameState.ActionCount}</h1>
        {this.hand(gameState, player, playerActions)}
        <br />
        {playerActions.discard ? "Select Card to Discard" : null}
        {this.cureButton(playerActions)}
        {this.getButton(playerActions)}
        {this.researchButton(playerActions)}
        {this.moveButton(playerActions)}
        <br />
        {this.epidemicButton(playerActions)}
        {this.drawButton(playerActions)}
        {this.infectButton(playerActions)}
        {this.outbreakButton(playerActions)}
      </div>
    )
  }
  render() {
    return (
      <div className="App">
        <div style={{ height: "50vh" }}>
          <ReactFlow
            elements={this.state.elements}
            nodesDraggable={false}
            nodesConnectable={false}
            elementsSelectable={true}
            preventScrolling={true}
            paneMoveable={false}
            zoomOnScroll={false}
            zoomOnPinch={false}
            zoomOnDoubleClick={false}
            onLoad={this.onLoad}
          />
        </div>
        <div className="Actions">
          {this.actions()}
        </div>
      </div>
    )
  }
}

export default App;