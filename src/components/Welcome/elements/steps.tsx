import React from 'react'
import { render } from 'react-dom'
import styled from 'styled-components'

type stepsPros = {
   step?: number,
}

export default function Steps({ step }: stepsPros) {
   let activeWelcome = false
   let activeLocation = false
   let activeKeyboard = false
   let activePartitions = false
   let activeUsers = false
   let activeNetwork = false
   let activeSummary = false
   let activeInstall = false
   let activeFinish = false

   if (step === 1) {
      activeWelcome = true
   } else if (step === 2) {
      activeLocation = true
   } else if (step === 3) {
      activeKeyboard = true
   } else if (step === 4) {
      activePartitions = true
   } else if (step === 5) {
      activeUsers = true
   } else if (step === 6) {
      activeNetwork = true
   } else if (step === 7) {
      activeSummary = true
   } else if (step === 8) {
      activeInstall = true
   } else if (step === 9) {
      activeFinish = true
   }

   return (
      <>
         <div >
            <WelcomeTab active={activeWelcome}/>
            <LocationTab active={activeLocation}/>
            <KeyboardTab active={activeKeyboard}/>
            <PartitionTab active={activePartitions}/>
            <UsersTab active={activeUsers} />
            <NetworkTab active={activeNetwork} />
            <SummaryTab active={activeSummary} />
            <InstallTab active={activeInstall} />
            <FinishTab active={activeFinish} />
         </div>
      </>
   )
}


type elementType = {
   active?: boolean
}

const Tab = styled.div`
width: 15vh;
background: ${props => props.activated ? "black" : "white"};
color: ${props => props.activated ? "white" : "black"};
`

function activateTab( label = '', active = false ): JSX.Element {
   let elem = <Tab >{label}</Tab>
   if (active) {
      elem = <Tab activated >{label}</Tab>
   }
   return elem
}


function WelcomeTab({ active = false }): JSX.Element {
   return activateTab('Welcome', active)
}

function LocationTab({ active = false }): JSX.Element {
   return activateTab('Location', active)
}

function KeyboardTab({ active = false }): JSX.Element {
   return activateTab('Keyboard', active)
}

function PartitionTab({ active = false }): JSX.Element {
   return activateTab('Partitions', active)
}

function UsersTab({ active = false }): JSX.Element {
   return activateTab('Users', active)
}

function NetworkTab({ active = false }): JSX.Element {
   return activateTab('Network', active)
}

function SummaryTab({ active = false }): JSX.Element {
   return activateTab('Summary', active)
}

function InstallTab({ active = false }): JSX.Element {
   return activateTab('Install', active)
}

function FinishTab({ active = false }): JSX.Element {
   return activateTab('Install', active)
}
