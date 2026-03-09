# Full-Stack Engineer Interview Transcript
**Interviewer:** John (Senior Engineering Manager)
**Candidate:** Suraj
**Position:** Full-Stack Engineer (Node.js + React.js)
**Duration:** ~75 minutes

**John:** Hey Suraj, welcome! Thanks for coming in — or well, logging on. Can you hear me okay?

**Suraj:** Yeah, loud and clear. Thanks for having me, John. Excited to be here.

**John:** Great. So, we'll have about 75 minutes together today. I want to cover a bit of your background, some technical concepts, a couple of coding problems, and then leave time at the end for your questions. Sound good?

**Suraj:** Sounds perfect.

**John:** Awesome. Let's start easy — tell me a little about yourself and your current role.

**Suraj:** Sure. So I'm currently a software engineer at a mid-sized fintech startup. I've been there for about three years. My role is technically full-stack, but I'd say I lean more toward the backend — probably 60-40 backend to frontend. I work primarily in Node.js on the backend, and we use React on the frontend. Before that, I was at a smaller agency for about two years where I did a lot of greenfield project work, which was a great learning experience but also kind of chaotic.

**John:** Ha, I can imagine. What does your current tech stack look like day to day?

**Suraj:** So on the backend, it's Node.js with Express, PostgreSQL as the primary database, and we use Redis for caching and session management. We also have some microservices talking to each other over RabbitMQ. On the frontend it's React with TypeScript, and we use Redux Toolkit for state management. We're on AWS — EC2, RDS, S3, and we've been slowly migrating some things to Lambda.

**John:** Nice. That's a pretty mature stack. What's the team size like?

**Suraj:** About 12 engineers total — three dedicated frontend, five backend including myself, and the rest are kind of mixed like me. We work in two-week sprints, pretty standard Agile setup.

**John:** Got it. And what made you start looking around?

**Suraj:** Honestly, it's been a great experience, but I feel like I've reached a ceiling in terms of the technical challenges. The systems are fairly stable now, and a lot of my work has become more maintenance-oriented. I want to be in an environment where I'm still being pushed, building at scale, solving harder problems. When I saw this role and the scale you're operating at, it immediately stood out.

**John:** That's fair. We definitely have no shortage of hard problems here. Okay, let's get into the technical stuff.

**John:** Let's start with JavaScript fundamentals. Can you explain the event loop to me? Pretend I'm a junior engineer who's just started.

**Suraj:** Sure, yeah. So JavaScript is single-threaded, meaning it can only execute one piece of code at a time. The event loop is the mechanism that allows it to handle asynchronous operations — like network requests, file reads, timers — without blocking the main thread.

Here's the mental model I use: imagine a call stack, a callback queue, and the event loop sitting between them. When your code runs, function calls get pushed onto the call stack and executed. When the stack is empty, the event loop checks the callback queue and pushes any pending callbacks onto the stack to be executed.

But it's a bit more nuanced than that. There's actually a microtask queue and a macrotask queue. Promises and `queueMicrotask` go into the microtask queue, while things like `setTimeout` and `setInterval` go into the macrotask queue. The key thing is that microtasks always get processed before the next macrotask. So if you have a resolved promise and a setTimeout both waiting, the promise callback fires first — even if the setTimeout was set to zero milliseconds.

**John:** Good. What would be the output of this:

```javascript
console.log('1');
setTimeout(() => console.log('2'), 0);
Promise.resolve().then(() => console.log('3'));
console.log('4');
```

**Suraj:** Okay, so — '1' logs immediately because it's synchronous. '4' logs next, same reason. Then the promise `.then` callback is in the microtask queue, so '3' logs before '2'. And finally '2' logs from the setTimeout macrotask. So: 1, 4, 3, 2.

**John:** Exactly right. Okay, next one — explain closures and give me a practical use case in a Node.js context.

**Suraj:** A closure is when a function retains access to variables from its outer lexical scope even after that outer function has returned. The inner function "closes over" those variables.

A practical Node.js example — rate limiting middleware. You can write a function that takes a limit and a time window as arguments, and inside it maintains a map of IP addresses to request counts. You return a middleware function that has closure over that map. Every request goes through the middleware, which updates the map and enforces the limit. The map persists across requests because of the closure — you don't need to store it externally for simple cases.

Another common one is creating private state in a module. Before ES modules and classes were widespread, closures were the standard way to encapsulate data.

**John:** Nice. Alright, let's talk about `async/await` versus Promises. Are there scenarios where you'd prefer one over the other?

**Suraj:** In most application code, I prefer `async/await` for readability. It reads more like synchronous code and makes error handling with try/catch feel natural. But there are cases where raw Promises are better.

If you need to run multiple independent async operations concurrently, `Promise.all` is much cleaner than trying to do that with multiple awaits. Like:

```javascript
const [user, orders] = await Promise.all([
  fetchUser(id),
  fetchOrders(id)
]);
```

That runs both in parallel. If you awaited them sequentially, you'd be adding latency unnecessarily.

Also, `Promise.allSettled` is useful when you want all operations to complete regardless of individual failures — like when you're making several third-party API calls and you don't want one failure to abort everything.

And sometimes for streaming or complex chaining, the Promise API feels more expressive. But 90% of the time, `async/await` is what I reach for.

**John:** What about error handling? What are some pitfalls with async/await?

**Suraj:** The big one is forgetting to handle rejections. If you `await` something and it throws and you don't have a try/catch, you get an unhandled promise rejection. In older versions of Node that was just a warning, but in newer versions it terminates the process — which is the right behavior, honestly, but it can catch people off guard.

Another one is accidentally running things sequentially when they could be parallel, which I touched on. And there's the issue of error swallowing — if you have a catch block that doesn't re-throw or properly handle the error, your code silently fails, which can be a nightmare to debug.

In Express specifically, one common gotcha is that if you have an async route handler that throws, Express won't automatically catch it and pass it to the error middleware unless you explicitly wrap it or use a library like `express-async-errors`. I've seen that bite junior engineers a lot.

**John:** That's a good point. We actually ran into that exact issue here early on. Okay, Node.js architecture — when would you use clustering or worker threads?

**Suraj:** They solve different problems. Clustering is for maximizing CPU utilization for I/O-bound workloads. Since Node is single-threaded, a single process only ever uses one CPU core. Clustering lets you fork multiple processes — typically one per core — each running its own event loop and handling requests independently. The `cluster` module or something like PM2 handles this. It's great for HTTP servers under high load.

Worker threads are for CPU-bound tasks within a single Node process — things like heavy computation, image processing, parsing large files. Before worker threads, doing CPU-heavy work blocked the event loop and made your server unresponsive. With worker threads, you offload that to a separate thread while the main event loop stays free. They share memory more efficiently than clusters through SharedArrayBuffer, which can be powerful but also requires careful synchronization.

At my current job, we use clustering for our main API servers via PM2, and we use a worker thread pool for a report generation feature that involves processing large CSV exports.

**John:** Let's move to React. Walk me through the component lifecycle with hooks — how do you think about `useEffect`?

**Suraj:** `useEffect` is how you synchronize your component with external systems — side effects like data fetching, subscriptions, DOM manipulation, timers. The mental model I've settled on is: don't think of it as lifecycle methods translated to hooks. Think of it as "what does this component need to stay in sync with?"

The dependency array controls when the effect re-runs. Empty array means it runs once after the initial render — similar to `componentDidMount`. Dependencies listed means it re-runs whenever those values change. No array at all means it runs after every render, which is usually not what you want.

The cleanup function is important — you return a function from `useEffect` that runs before the next effect fires or when the component unmounts. This is where you clear timers, cancel subscriptions, abort fetch requests. Forgetting cleanup is a common source of memory leaks and the "Can't perform a React state update on an unmounted component" warning — though technically React 18 made that warning less strict.

**John:** What are some common mistakes you see with hooks in general?

**Suraj:** A few. Stale closures are probably the most insidious — where a callback inside an effect captures an old value of state or props because the effect only ran once and closed over the initial values. The fix is usually adding the dependency to the array or using `useRef` to hold a mutable reference.

Violating the rules of hooks is another one — calling hooks conditionally or inside loops. That breaks the order guarantee React depends on. The ESLint plugin for hooks catches most of these.

Overusing `useEffect` for things that should just be derived values is common too. If something can be computed during render from existing state and props, you don't need an effect. The React team has actually been pretty vocal about this in their docs recently.

And `useCallback` and `useMemo` misuse — people reach for them reflexively thinking they're optimizing, but they have overhead of their own. The rule of thumb I follow is: only memoize when you can measure a performance problem or when you're passing callbacks to deeply memoized child components where referential stability matters.

**John:** Good. Let's talk about state management. When do you reach for Redux versus React context versus something like Zustand?

**Suraj:** Context is great for low-frequency updates — theme, locale, auth state, current user. Things that don't change often and need to be accessible throughout the tree. The issue with Context is that any component consuming the context re-renders when it changes, so you don't want to put high-frequency state in there without memoization strategies.

Redux — or specifically Redux Toolkit in modern apps — makes sense when you have complex, interconnected global state that multiple parts of the app depend on, you need powerful dev tooling like time-travel debugging, or you have complex async flows with RTK Query or Redux Saga. It adds boilerplate, but that structure pays dividends on large teams.

Zustand is somewhere in the middle — minimal boilerplate, no provider needed, really clean API. I've used it for medium-complexity apps where Redux felt like overkill. It handles subscriptions intelligently, so components only re-render when the specific slice of state they use changes.

The honest answer is: I try to start with local state, elevate to Context if needed, and only bring in a full state management library when the complexity justifies it.

**John:** How do you approach performance optimization in React?

**Suraj:** First, measure. React DevTools Profiler is the starting point — identify which components are re-rendering unnecessarily and how expensive those renders are. Don't optimize blindly.

For preventing unnecessary re-renders: `React.memo` for components, `useMemo` for expensive derived values, `useCallback` for stable function references passed to memoized children.

Code splitting with `React.lazy` and `Suspense` is huge for initial load performance — dynamically import components that aren't needed upfront. Route-based splitting is the low-hanging fruit.

List virtualization — if you're rendering hundreds or thousands of items, react-window or react-virtual are essential. Rendering ten thousand DOM nodes kills performance.

And beyond React-specific stuff: image optimization, avoiding layout thrashing, keeping bundle size in check with tree shaking, lazy loading images. The fundamentals still matter.

**John:** Solid. One more React question — what's your understanding of React 18's concurrent features?

**Suraj:** React 18 introduced the concurrent renderer, which allows React to prepare multiple versions of the UI simultaneously and interrupt renders to handle more urgent updates. The key idea is prioritized rendering.

`useTransition` and `startTransition` let you mark state updates as non-urgent — like filtering a large list. React can start that render, pause it if a more urgent update comes in — like a user typing — and resume later. This keeps the UI responsive.

`useDeferredValue` is similar but for deferring a value rather than a transition. You can use it to show the previous content while new content is being prepared.

Suspense got more powerful with concurrent mode — you can now use it for data fetching with frameworks that support it, not just lazy-loaded components. And Suspense boundaries can be nested to give you fine-grained loading states.

Automatic batching is another React 18 change — state updates in async callbacks, timeouts, and native event handlers are now batched automatically, whereas before only updates in React event handlers were batched. That's a nice performance improvement without any code changes.

**John:** Okay, let's do a mini system design. You're building a notification service — users can receive notifications via email, SMS, and in-app push. It needs to handle high volume — say a million notifications per day. Walk me through how you'd design this.

**Suraj:** Sure. Let me think through this out loud.

First, I'd separate the concerns into a few pieces: the notification API that accepts notification requests, a message queue for durability and decoupling, channel-specific workers, and a delivery tracking store.

The API layer would be a simple Node.js service — clients submit a notification request with a payload: recipient, channel, type, content. The API validates the request, persists a notification record to the database with a "pending" status, and publishes a message to a queue — something like RabbitMQ or SQS. Responding quickly to the caller is important; we don't want them waiting for delivery.

The queue is the backbone. It decouples the API from the workers and gives us durability — if a worker crashes, the message isn't lost. With a million notifications a day, that's around twelve per second on average, but you'd want to handle spike capacity — maybe ten to fifty times that during peak. SQS handles that easily.

For workers, I'd have separate consumer services for each channel — an email worker, an SMS worker, an in-app worker. Each pulls from the queue, calls the appropriate third-party provider — SendGrid for email, Twilio for SMS, your own WebSocket layer or Firebase for push — and updates the notification status in the database. They'd also handle retries with exponential backoff for transient failures and dead-letter queues for messages that consistently fail.

The database — I'd probably use PostgreSQL for the notification records, with indexes on user ID and status for querying. Redis for deduplication — maintain a short TTL cache of recently sent notification IDs to catch duplicates from at-least-once delivery semantics. And if you need analytics or a queryable log at scale, something like ClickHouse or BigQuery for the delivery history.

For in-app notifications specifically, you need a real-time component. WebSockets via Socket.io or Server-Sent Events. The in-app worker publishes to a Redis pub/sub channel, and your WebSocket servers are subscribed. When a user is connected, they receive the notification in real time. If they're offline, it's stored in the database and served when they reconnect.

A few other considerations: rate limiting per user or per channel to avoid spamming, user preference management — respecting do-not-disturb settings, channel opt-outs — and observability: metrics on delivery rates, latency, failure rates by channel.

**John:** Good answer. How would you handle idempotency in this system?

**Suraj:** Critical point. With at-least-once delivery from the queue, workers could process the same message more than once — due to network hiccups or worker restarts before acknowledgment. You absolutely don't want to send someone two identical emails.

The pattern I'd use: include a unique idempotency key in each notification record — either a UUID generated at creation time or a content-based hash. Before the worker sends via the third-party API, it checks Redis with that key. If the key exists with a "sent" status, it skips and acknowledges the message. If not, it proceeds, and on success it sets the key in Redis with a TTL long enough to cover your retry window — say 24 to 48 hours.

Some third-party APIs like Stripe and some email providers support idempotency keys natively, so you can also pass your key through to them. Belt and suspenders.

**John:** Makes sense. What about database design — let's say we need to store user permissions and roles for a multi-tenant SaaS application. How do you approach that schema?

**Suraj:** Classic RBAC — Role-Based Access Control — is the starting point. The core entities are: Users, Roles, Permissions, and then the joining tables.

So you'd have a `tenants` table since it's multi-tenant. Users belong to a tenant. Roles are scoped to a tenant too, because Tenant A's "admin" role might have different permissions than Tenant B's. Then `permissions` is a relatively static table of named capabilities — things like `reports:read`, `users:write`, `billing:manage`.

The join tables: `user_roles` maps users to roles within a tenant, and `role_permissions` maps roles to permissions. When you need to check if a user can do something, you query: get the user's roles for this tenant, get all permissions for those roles, check if the required permission is in the set.

For performance, I'd cache the resolved permission set per user in Redis on login and invalidate on role changes. Hitting the database for every authorization check doesn't scale.

If you need more granularity — like row-level permissions, where a user can edit their own records but not others' — you layer on an ABAC component, Attribute-Based Access Control. That gets into policy engines like Casbin or OPA, depending on complexity.

I'd also add an audit log table for permission changes — in a SaaS product, especially fintech or anything compliance-heavy, you need to know who granted what, when.

**John:** Okay, let's do a bit of live coding. Don't worry about it being perfect — I'm more interested in how you think. Can you share your screen?

**Suraj:** Sure, one sec. Okay, you should see my editor.

**John:** Perfect. First problem: implement a function that takes an array of async functions and executes them with a concurrency limit. So like a promise pool.

**Suraj:** Okay, classic concurrency control. Let me think through this...

So the signature would be something like `asyncPool(limit, tasks)` where `tasks` is an array of functions that return promises. We want at most `limit` running simultaneously.

```javascript
async function asyncPool(limit, tasks) {
  const results = [];
  const executing = new Set();

  for (const task of tasks) {
    const promise = Promise.resolve().then(() => task());
    results.push(promise);
    executing.add(promise);

    const cleanup = () => executing.delete(promise);
    promise.then(cleanup, cleanup);

    if (executing.size >= limit) {
      await Promise.race(executing);
    }
  }

  return Promise.all(results);
}
```

So — we maintain a Set of currently executing promises. For each task, we start it and add its promise to both the results array and the executing set. We attach a cleanup to remove it from executing when it settles. If we've hit the concurrency limit, we `await Promise.race(executing)`, which waits until at least one finishes before we add the next. At the end, `Promise.all(results)` waits for everything and returns all results.

**John:** Nice. What if one of the tasks throws — how does this behave?

**Suraj:** Right now, `Promise.all` at the end would reject with the first error. If you want all tasks to complete regardless of individual failures — fault-tolerant behavior — you'd change `Promise.all` to `Promise.allSettled` and you'd also want the cleanup to handle rejections, which it already does since I pass the same `cleanup` as both the success and error handler.

Actually, you'd probably also want to wrap each task call in a try/catch or `.catch` so a rejection doesn't propagate through `Promise.race` unexpectedly. Let me adjust...

```javascript
const promise = Promise.resolve().then(() => task()).catch(err => ({ error: err }));
```

Now failures don't abort the pool, and you can check for the `error` property in results downstream.

**John:** Great. Second problem — and this is more of a React one: implement a custom `useFetch` hook that handles loading, error, and data states, and also supports cancellation if the component unmounts.

**Suraj:** Sure. AbortController is the key here.

```javascript
import { useState, useEffect, useRef } from 'react';

function useFetch(url, options = {}) {
  const [state, setState] = useState({
    data: null,
    loading: true,
    error: null,
  });
  
  const optionsRef = useRef(options);

  useEffect(() => {
    const controller = new AbortController();
    const signal = controller.signal;

    setState({ data: null, loading: true, error: null });

    fetch(url, { ...optionsRef.current, signal })
      .then(res => {
        if (!res.ok) throw new Error(`HTTP error: ${res.status}`);
        return res.json();
      })
      .then(data => {
        if (!signal.aborted) {
          setState({ data, loading: false, error: null });
        }
      })
      .catch(err => {
        if (err.name !== 'AbortError') {
          setState({ data: null, loading: false, error: err });
        }
      });

    return () => controller.abort();
  }, [url]);

  return state;
}
```

So — we create an `AbortController` per effect invocation. The signal gets passed into fetch. The cleanup function calls `controller.abort()`, which cancels the in-flight request on unmount or if the URL changes before the request completes. We check `signal.aborted` before updating state as an extra guard. And we filter out `AbortError` because that's not a real error from the consumer's perspective — it's expected cleanup behavior.

I'm storing options in a ref so they don't need to be in the dependency array — passing an inline object as the options argument on every render would cause an infinite loop since objects are compared by reference.

**John:** Smart. What would you add to make this production-ready?

**Suraj:** A few things. Caching — for the same URL, you don't want to re-fetch on every component mount. You could maintain a simple in-memory cache or integrate with something like SWR or React Query, which handle all of this plus revalidation strategies.

Retry logic for transient failures. Exponential backoff with a max retry count.

Stale-while-revalidate — show cached data immediately while fetching fresh data in the background. Much better UX than blanking out to a loading state every time.

Deduplication — if ten components mount at the same time and all call this hook with the same URL, you don't want ten concurrent requests. Deduplicate in-flight requests.

And honestly, all of this exists in React Query or SWR already, which is why I'd recommend those over rolling your own in most production applications. But writing it from scratch like this is a great exercise for understanding what's going on under the hood.

**John:** Totally agree. Okay, we're almost at time. You've done really well today — solid answers across the board. Before we wrap up, do you have any questions for me?

**Suraj:** Yeah, a few if that's okay. First — what does the engineering culture around code review look like? I find that tells you a lot about a team.

**John:** It's pretty healthy here. We aim for small PRs and fast turnarounds — ideally same day, definitely within 24 hours. Reviews are expected to be constructive, not just nitpicky. We have automated checks for lint, tests, and coverage that gate merges, so human reviews can focus on design and logic rather than formatting. There's a strong culture of actually explaining the "why" in comments, not just the "what."

**Suraj:** That's great to hear. Second — what are the biggest technical challenges the team is working through right now?

**John:** Honestly, scale. Our user base grew about four times in the last 18 months, and some of our earlier architectural decisions are starting to show stress. We're working on improving observability across our services, breaking apart a couple of monoliths, and figuring out our data layer as query volumes grow. Lots of meaty problems.

**Suraj:** That's exactly the kind of environment I'm looking for. Last one — how does the team approach technical debt?

**John:** We try to be deliberate about it. We have a concept of "tech debt Fridays" — not mandatory, but engineers are encouraged to spend time on cleanup, refactoring, improving test coverage. We also have a running backlog of debt items that get prioritized alongside feature work each sprint. We don't let it just pile up silently — it's tracked, estimated, and treated like real work.

**Suraj:** Love that. That kind of structure around it makes a big difference.

**John:** Agreed. Alright, Suraj — this was a great conversation. I'll connect with you via the recruiter this week. Really enjoyed the discussion.

**Suraj:** Likewise, John. Thanks so much for the time — really appreciated the depth of the conversation.

**John:** Of course. Take care!

**Suraj:** You too. Bye!
