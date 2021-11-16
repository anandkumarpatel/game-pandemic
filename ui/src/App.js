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

  onLoad(reactFlowInstance) {
    const elements = []
    return fetch("http://localhost:8080/state")
      .then(res => res.json())
      .then((input) => {

        const state = input.State
        console.log(input.State)
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
                  width: "1px",
                  height: "1px",
                  background: vName,
                },
                position: {
                  x: cityX + 15 * xMap[vName] + 15 * i,
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
              x: cityX + 100,
              y: cityY - 10 + 10 * pPos[pName],
            },
          })

        })


        this.setState({ elements, game: input }, () => {
          console.log("state", this.state)
        })
      })
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

      >
        {playerActions.move.map((city) => {
          return <Dropdown.Item key={city} eventKey={`move-${city}`}>{city}</Dropdown.Item>
        })}
      </DropdownButton>
    )
  }

  cureButton(playerActions) {
    return (
      <DropdownButton
        size="lg"
        drop="down"
        title="cure"
        disabled={!playerActions.cure}
        as={ButtonGroup}
      >
        {playerActions.cure && playerActions.cure.map((city) => {
          return <Dropdown.Item key={city} eventKey={`cure-${city}`}>{city}</Dropdown.Item>
        })}
      </DropdownButton>
    )
  }

  infectButton(playerActions) {
    return (
      <Button
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
    return (
      <Button
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
    return (
      <Button
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
    return (
      <Button
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

  reseachButton(playerActions) {
    return (
      <DropdownButton
        size="lg"
        drop="down"
        title="reseach"
        disabled={!playerActions.reseach}
        as={ButtonGroup}
      >
        {playerActions.reseach && playerActions.reseach.map((city) => {
          return <Dropdown.Item key={city} eventKey={`research-${city}`}>{city}</Dropdown.Item>
        })}
      </DropdownButton>
    )
  }

  getButton(playerActions) {
    return (
      <DropdownButton
        size="lg"
        drop="down"
        title="get"
        disabled={!playerActions.get}
        as={ButtonGroup}
      >
        {playerActions.get && playerActions.get.map((city) => {
          return <Dropdown.Item key={city} eventKey={`get-${city}`}>{city}</Dropdown.Item>
        })}
      </DropdownButton>
    )
  }

  actions() {
    if (!this.state.game) return null
    const player = this.state.game.State.Players[this.state.game.State.CurrentPlayerN]
    const playerActions = this.state.game.Actions[player.Name]

    return (
      <div>
        <h1>Player {player.Name}</h1>
        {player.Hand.Cards.map((card) => {
          console.log(playerActions.flyTo.includes(card.Name))
          return (
            <DropdownButton
              variant={card.VirusType.toLowerCase()}
              size="lg"
              drop="up"
              key={card.Name}
              as={ButtonGroup}
              title={card.Name}
            >
              {playerActions.flyTo && playerActions.flyTo.includes(card.Name) ? <Dropdown.Item key={`flyTo-${card.Name}`} eventKey={`flyTo-${card.Name}`}>Fly</Dropdown.Item> : null}
              {playerActions.build && playerActions.build.includes(card.Name) ? <Dropdown.Item key={`build-${card.Name}`} eventKey={`build-${card.Name}`}>Build</Dropdown.Item> : null}
              {playerActions.flyAnywhere && playerActions.flyAnywhere.includes(card.Name) ? <Dropdown.Item key={`flyAnywhere-${card.Name}`} eventKey={`flyAnywhere-${card.Name}`}>Build</Dropdown.Item> : null}
              {playerActions.discard && playerActions.discard.includes(card.Name) ? <Dropdown.Item key={`discard-${card.Name}`} eventKey={`discard-${card.Name}`}>Build</Dropdown.Item> : null}
              {playerActions.giveCard && playerActions.giveCard.includes(card.Name) ? <Dropdown.Item key={`giveCard-${card.Name}`} eventKey={`giveCard-${card.Name}`}>Give</Dropdown.Item> : null}
            </DropdownButton>
          )
        })}
        <br />
        {this.moveButton(playerActions)}
        {this.cureButton(playerActions)}
        {this.getButton(playerActions)}
        {this.reseachButton(playerActions)}
        {this.epidemicButton(playerActions)}
        <br />
        {this.drawButton(playerActions)}
        <br />
        {this.infectButton(playerActions)}
        <br />
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
