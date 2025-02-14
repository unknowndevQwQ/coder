import Button from "@material-ui/core/Button"
import Link from "@material-ui/core/Link"
import Popover from "@material-ui/core/Popover"
import { makeStyles } from "@material-ui/core/styles"
import TextField from "@material-ui/core/TextField"
import OpenInNewOutlined from "@material-ui/icons/OpenInNewOutlined"
import { ChooseOne, Cond } from "components/Conditionals/ChooseOne"
import { Stack } from "components/Stack/Stack"
import { useRef, useState, Fragment } from "react"
import { colors } from "theme/colors"
import { CodeExample } from "../CodeExample/CodeExample"
import {
  HelpTooltipLink,
  HelpTooltipLinksGroup,
  HelpTooltipText,
  HelpTooltipTitle,
} from "../Tooltips/HelpTooltip"
import { Maybe } from "components/Conditionals/Maybe"
import { useMachine } from "@xstate/react"
import { portForwardMachine } from "xServices/portForward/portForwardXService"

export interface PortForwardButtonProps {
  host: string
  username: string
  workspaceName: string
  agentName: string
  agentId: string
}

const EnabledView: React.FC<PortForwardButtonProps> = (props) => {
  const { host, workspaceName, agentName, agentId, username } = props
  const styles = useStyles()
  const [port, setPort] = useState("3000")
  const { location } = window
  const urlExample = `${location.protocol}//${port}--${agentName}--${workspaceName}--${username}.${host}`
  const [state] = useMachine(portForwardMachine, {
    context: { agentId: agentId },
  })
  const ports = state.context.listeningPorts?.ports

  return (
    <>
      <HelpTooltipText>
        Access ports running on the agent with the{" "}
        <strong>port, agent name, workspace name</strong> and{" "}
        <strong>your username</strong> URL schema, as shown below.
      </HelpTooltipText>

      <CodeExample code={urlExample} className={styles.code} />

      <HelpTooltipText>
        Use the form to open applications in a new tab.
      </HelpTooltipText>

      <Stack
        direction="row"
        spacing={1}
        alignItems="center"
        className={styles.form}
      >
        <TextField
          label="Port"
          type="number"
          value={port}
          className={styles.portField}
          onChange={(e) => {
            setPort(e.currentTarget.value)
          }}
        />
        <Link
          underline="none"
          href={urlExample}
          target="_blank"
          rel="noreferrer"
          className={styles.openUrlButton}
        >
          <Button>Open URL</Button>
        </Link>
      </Stack>

      <Maybe condition={Boolean(ports && ports.length > 0)}>
        <HelpTooltipText>
          {ports &&
            ports.map((p, i) => {
              const url = `${location.protocol}//${p.port}--${agentName}--${workspaceName}--${username}.${host}`
              let label = `${p.port}`
              if (p.process_name) {
                label = `${p.process_name} - ${p.port}`
              }

              return (
                <Fragment key={i}>
                  {i > 0 && <span style={{ margin: "0 0.6em" }}>&middot;</span>}
                  <Link href={url} target="_blank" rel="noreferrer">
                    {label}
                  </Link>
                </Fragment>
              )
            })}
        </HelpTooltipText>
      </Maybe>

      <HelpTooltipLinksGroup>
        <HelpTooltipLink href="https://coder.com/docs/coder-oss/latest/networking/port-forwarding#dashboard">
          Learn more about web port forwarding
        </HelpTooltipLink>
      </HelpTooltipLinksGroup>
    </>
  )
}

const DisabledView: React.FC<PortForwardButtonProps> = ({
  workspaceName,
  agentName,
}) => {
  const cliExample = `coder port-forward ${workspaceName}.${agentName} --tcp 3000`
  const styles = useStyles()

  return (
    <>
      <HelpTooltipText>
        <strong>
          Your deployment does not have web port forwarding enabled.
        </strong>{" "}
        See the docs for more details.
      </HelpTooltipText>

      <HelpTooltipText>
        You can use the Coder CLI to forward ports from your workspace to your
        local machine, as shown below.
      </HelpTooltipText>

      <CodeExample code={cliExample} className={styles.code} />

      <HelpTooltipLinksGroup>
        <HelpTooltipLink href="https://coder.com/docs/coder-oss/latest/networking/port-forwarding#dashboard">
          Learn more about web port forwarding
        </HelpTooltipLink>
      </HelpTooltipLinksGroup>
    </>
  )
}

export const PortForwardButton: React.FC<PortForwardButtonProps> = (props) => {
  const { host } = props
  const anchorRef = useRef<HTMLButtonElement>(null)
  const [isOpen, setIsOpen] = useState(false)
  const id = isOpen ? "schedule-popover" : undefined
  const styles = useStyles()

  const onClose = () => {
    setIsOpen(false)
  }

  return (
    <>
      <Button
        startIcon={<OpenInNewOutlined />}
        size="small"
        ref={anchorRef}
        onClick={() => {
          setIsOpen(true)
        }}
      >
        Port forward
      </Button>
      <Popover
        classes={{ paper: styles.popoverPaper }}
        id={id}
        open={isOpen}
        anchorEl={anchorRef.current}
        onClose={onClose}
        anchorOrigin={{
          vertical: "bottom",
          horizontal: "left",
        }}
        transformOrigin={{
          vertical: "top",
          horizontal: "left",
        }}
      >
        <HelpTooltipTitle>Port forward</HelpTooltipTitle>
        <ChooseOne>
          <Cond condition={host !== ""}>
            <EnabledView {...props} />
          </Cond>
          <Cond>
            <DisabledView {...props} />
          </Cond>
        </ChooseOne>
      </Popover>
    </>
  )
}

const useStyles = makeStyles((theme) => ({
  popoverPaper: {
    padding: `${theme.spacing(2.5)}px ${theme.spacing(3.5)}px ${theme.spacing(
      3.5,
    )}px`,
    width: theme.spacing(52),
    color: theme.palette.text.secondary,
    marginTop: theme.spacing(0.25),
  },

  openUrlButton: {
    flexShrink: 0,
  },

  portField: {
    // The default border don't contrast well with the popover
    "& .MuiOutlinedInput-root .MuiOutlinedInput-notchedOutline": {
      borderColor: colors.gray[10],
    },
  },

  code: {
    margin: theme.spacing(2, 0),
  },

  form: {
    margin: theme.spacing(1.5, 0, 0),
  },
}))
