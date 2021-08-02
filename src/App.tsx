import { GlobalStyle } from './styles/GlobalStyle'

import { Greetings } from './components/Greetings'
import Welcome from './components/Welcome/welcome'

export function App() {
  return (
    <>
      <GlobalStyle />
      <Welcome />
    </>
  )
}