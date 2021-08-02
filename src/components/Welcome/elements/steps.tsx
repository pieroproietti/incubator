import React from 'react'
import { render } from 'react-dom'


type stepsPros = {
   step?: number,
}

export default function Steps({ step = 1 }: stepsPros) {
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
            <WelcomeTab active={activeWelcome} />
            <LocationTab active={activeLocation} />
            <KeyboardTab active={activeKeyboard} />
            <PartitionTab active={activePartitions} />
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

function WelcomeTab({ active = false }): JSX.Element {
   let backgroundColor = 'white'
   let color = 'black'
   if (active) {
      backgroundColor = 'black'
      color = 'white'
   }
   return <div color={color} background-color={backgroundColor}> Welcome    </div>
}

function LocationTab({ active = false }): JSX.Element {
   let backgroundColor = 'white'
   let color = 'black'
   if (active) {
      backgroundColor = 'black'
      color = 'white'
   }
   return <div color={color} background-color={backgroundColor}> Location   </div>
}

function KeyboardTab({ active = false }): JSX.Element {
   let backgroundColor = 'white'
   let color = 'black'
   if (active) {
      backgroundColor = 'black'
      color = 'white'
   }
   return <div color={color} background-color={backgroundColor}> Keyboard   </div>
}


function PartitionTab({ active = false }): JSX.Element {
   let backgroundColor = 'white'
   let color = 'black'
   if (active) {
      backgroundColor = 'black'
      color = 'white'
   }
   return <div color={color} background-color={backgroundColor}> Partitions </div>
}

function UsersTab({ active = false }): JSX.Element {
   let backgroundColor = 'white'
   let color = 'black'
   if (active) {
      backgroundColor = 'black'
      color = 'white'
   }
   return <div color={color} background-color={backgroundColor}> Users      </div>
}

function NetworkTab({ active = false }): JSX.Element {
   let backgroundColor = 'white'
   let color = 'black'
   if (active) {
      backgroundColor = 'black'
      color = 'white'
   }
   return <div color={color} background-color={backgroundColor}> Network    </div>
}

function SummaryTab({ active = false }): JSX.Element {
   let backgroundColor = 'white'
   let color = 'black'
   if (active) {
      backgroundColor = 'black'
      color = 'white'
   }
   return <div color={color} background-color={backgroundColor}> Summary    </div>
}

function InstallTab({ active = false }): JSX.Element {
   let backgroundColor = 'white'
   let color = 'black'
   if (active) {
      backgroundColor = 'black'
      color = 'white'
   }
   return <div color={color} background-color={backgroundColor}> Install    </div>
}

function FinishTab({ active = false }): JSX.Element {
   let backgroundColor = 'white'
   let color = 'black'
   if (active) {
      backgroundColor = 'black'
      color = 'white'
   }
   return <div color={color} background-color={backgroundColor}> Finish     </div>
}
