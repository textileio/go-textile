import UIKit
import Mobile

class ViewController: UIViewController {

  var node: MobileNode?

  override func viewDidLoad() {
    super.viewDidLoad()
    let paths = NSSearchPathForDirectoriesInDomains(.libraryDirectory, .userDomainMask, true)
    let dataDir = paths[0]

    DispatchQueue.global().async {
      self.node = MobileNewTextile(dataDir)
      if let node = self.node {
        try! node.start()
      }
    }
  }

  override func didReceiveMemoryWarning() {
    super.didReceiveMemoryWarning()
    DispatchQueue.global().async {
      if let node = self.node {
        try! node.stop()
      }
    }
  }

}
