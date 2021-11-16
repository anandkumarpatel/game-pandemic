import './App.css';
import ReactFlow from 'react-flow-renderer';

import citypos from './citypos.json'
import React from 'react';


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
      .then((data) => {
        data.Cities.forEach((city) => {
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

        data.Players.forEach((player) => {
          const pName = player.Name
          const cityName = player.Location
          const cityX = citypos[cityName].x * 150
          const cityY = citypos[cityName].y * 100
          console.log("XX", {
            x: cityX + 100,
            y: cityY + 5 * pPos[pName],
          }, pPos, pName)
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


        this.setState({ elements, game: data }, () => {
          console.log("state", this.state)
        })
      })
      .then(() => setTimeout(() => {
        reactFlowInstance.fitView()
      }, 1))
      .catch(e => console.log("error", e))
  }

  actions() {
    if (!this.state.game) return null
    const player = this.state.game.Players[this.state.game.CurrentPlayerN]
    return (
      <div>
        <h1>Player {player.Name}</h1>
        {player.Hand.Cards.map((card) => {
          console.log("card.VirusType", card.VirusType)
          const className = `card background-${card.VirusType.toLowerCase()} text-${textColorMap[card.VirusType]}`
          return <div key={card.Name} class={className}>{card.Name}</div>
        })}
        <br/>
        <button >Cure</button>
        <button >Move</button>
        <button >Build</button>
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
