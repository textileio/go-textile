//  Created by react-native-create-bridge

#import "TextileIPFS.h"
#import <Mobile/Mobile.h>

// import RCTBridge
#if __has_include(<React/RCTBridge.h>)
#import <React/RCTBridge.h>
#elif __has_include(“RCTBridge.h”)
#import “RCTBridge.h”
#else
#import “React/RCTBridge.h” // Required when used as a Pod in a Swift project
#endif

// import RCTEventDispatcher
#if __has_include(<React/RCTEventDispatcher.h>)
#import <React/RCTEventDispatcher.h>
#elif __has_include(“RCTEventDispatcher.h”)
#import “RCTEventDispatcher.h”
#else
#import “React/RCTEventDispatcher.h” // Required when used as a Pod in a Swift project
#endif

@interface TextileIPFS()

@property (nonatomic, strong) MobileNode *node;

@end

@implementation TextileIPFS
@synthesize bridge = _bridge;

// Export a native module
// https://facebook.github.io/react-native/docs/native-modules-ios.html
RCT_EXPORT_MODULE();

// Export constants
// https://facebook.github.io/react-native/releases/next/docs/native-modules-ios.html#exporting-constants
- (NSDictionary *)constantsToExport
{
  return @{
           @"EXAMPLE": @"example"
         };
}

// Export methods to a native module
// https://facebook.github.io/react-native/docs/native-modules-ios.html

RCT_EXPORT_METHOD(createNodeWithDataDir:(NSString *)dataDir)
{
  [self _createNodeWithDataDir:dataDir];
}

RCT_REMAP_METHOD(startNode, startNodeWithResolver:(RCTPromiseResolveBlock)resolve rejecter:(RCTPromiseRejectBlock)reject)
{
  BOOL success = [self _startNode];
  if(success) {
    resolve(@YES);
  } else {
    NSError *error = [NSError errorWithDomain:@"ipfs" code:0 userInfo:nil];
    reject(@"failed_to_start_node", @"Failed to start node", error);
  }
}

RCT_REMAP_METHOD(peerId, peerIdWithResolver:(RCTPromiseResolveBlock)resolve rejecter:(RCTPromiseRejectBlock)reject)
{
  NSString *peerId = [self _peerId];
  if(peerId) {
    resolve(peerId);
  } else {
    NSError *error = [NSError errorWithDomain:@"ipfs" code:1 userInfo:nil];
    reject(@"nil_peer_id", @"Peer id is undefined", error);
  }
}

RCT_REMAP_METHOD(key, keyWithResolver:(RCTPromiseResolveBlock)resolve rejecter:(RCTPromiseRejectBlock)reject)
{
  NSString *key = [self _key];
  if(key) {
    resolve(key);
  } else {
    NSError *error = [NSError errorWithDomain:@"ipfs" code:2 userInfo:nil];
    reject(@"nil_key", @"Key is undefined", error);
  }
}

RCT_EXPORT_METHOD(addImageAtPath:(NSString *)path resolver:(RCTPromiseResolveBlock)resolve rejecter:(RCTPromiseRejectBlock)reject)
{
  NSDictionary *result = [self _addImageAtPath:path];
  if(result) {
    resolve(result);
  } else {
    NSError *error = [NSError errorWithDomain:@"ipfs" code:3 userInfo:nil];
    reject(@"nil_add_image_result", @"Add image result is undefined", error);
  }
}

RCT_EXPORT_METHOD(exampleMethod)
{
  [self emitMessageToRN:@"EXAMPLE_EVENT" :nil];
}

// List all your events here
// https://facebook.github.io/react-native/releases/next/docs/native-modules-ios.html#sending-events-to-javascript
- (NSArray<NSString *> *)supportedEvents
{
  return @[@"SampleEvent"];
}

#pragma mark - Private methods

- (void)_createNodeWithDataDir:(NSString *)dataDir {
  self.node = MobileNewTextile(dataDir);
}

- (BOOL)_startNode {
  NSError *error;
  BOOL success = [self.node start:&error];
  return success;
}

- (NSString *)_peerId {
  return @"somepeerid";
}

- (NSString *)_key {
  return @"thisissomekey";
}

- (NSDictionary *)_addImageAtPath:(NSString *)path {
  return @{
           @"hash" : @"somehash",
           @"previewImageData" : @"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAASwAAAEsCAYAAAB5fY51AAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAphgAAKYYBIuzfjAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAASPSURBVHic7daxDeVGFATBIT3ln6vW/OcoAFLGEQ1UOWvOWo13bft32z/7u85/r127du0+3r22/f7yKMD/cn/9AYCnBAvIECwgQ7CADMECMgQLyBAsIEOwgAzBAjIEC8gQLCBDsIAMwQIyBAvIECwgQ7CADMECMgQLyBAsIEOwgAzBAjIEC8gQLCBDsIAMwQIyBAvIECwg4952Ptg9du3atft214UFZFzbfl9/AuAJFxaQIVhAhmABGYIFZAgWkCFYQIZgARmCBWQIFpAhWECGYAEZggVkCBaQIVhAhmABGYIFZAgWkCFYQIZgARmCBWQIFpAhWECGYAEZggVkCBaQIVhAhmABGfe288HusWvXrt23uy4sIOPa9vv6EwBPuLCADMECMgQLyBAsIEOwgAzBAjIEC8gQLCBDsIAMwQIyBAvIECwgQ7CADMECMgQLyBAsIEOwgAzBAjIEC8gQLCBDsIAMwQIyBAvIECwgQ7CADMECMgQLyLi3nQ92j127du2+3XVhARnXtt/XnwB4woUFZAgWkCFYQIZgARmCBWQIFpAhWECGYAEZggVkCBaQIVhAhmABGYIFZAgWkCFYQIZgARmCBWQIFpAhWECGYAEZggVkCBaQIVhAhmABGYIFZAgWkCFYQMa97Xywe+zatWv37a4LC8i4tv2+/gTAEy4sIEOwgAzBAjIEC8gQLCBDsIAMwQIyBAvIECwgQ7CADMECMgQLyBAsIEOwgAzBAjIEC8gQLCBDsIAMwQIyBAvIECwgQ7CADMECMgQLyBAsIEOwgAzBAjLubeeD3WPXrl27b3ddWEDGte339ScAnnBhARmCBWQIFpAhWECGYAEZggVkCBaQIVhAhmABGYIFZAgWkCFYQIZgARmCBWQIFpAhWECGYAEZggVkCBaQIVhAhmABGYIFZAgWkCFYQIZgARmCBWQIFpBxbzsf7B67du3afbvrwgIyrm2/rz8B8IQLC8gQLCBDsIAMwQIyBAvIECwgQ7CADMECMgQLyBAsIEOwgAzBAjIEC8gQLCBDsIAMwQIyBAvIECwgQ7CADMECMgQLyBAsIEOwgAzBAjIEC8gQLCBDsICMe9v5YPfYtWvX7ttdFxaQcW37ff0JgCdcWECGYAEZggVkCBaQIVhAhmABGYIFZAgWkCFYQIZgARmCBWQIFpAhWECGYAEZggVkCBaQIVhAhmABGYIFZAgWkCFYQIZgARmCBWQIFpAhWECGYAEZggVk3NvOB7vHrl27dt/uurCAjGvb7+tPADzhwgIyBAvIECwgQ7CADMECMgQLyBAsIEOwgAzBAjIEC8gQLCBDsIAMwQIyBAvIECwgQ7CADMECMgQLyBAsIEOwgAzBAjIEC8gQLCBDsIAMwQIyBAvIECwg4952Ptg9du3atft214UFZFzbfl9/AuAJFxaQIVhAhmABGYIFZAgWkCFYQIZgARmCBWQIFpAhWECGYAEZggVkCBaQIVhAhmABGYIFZAgWkCFYQIZgARmCBWQIFpAhWECGYAEZggVkCBaQIVhAhmABGfe288HusWvXrt23u38AiEfYMqEQXI8AAAAASUVORK5CYII="
           };
}

// Implement methods that you want to export to the native module
- (void) emitMessageToRN: (NSString *)eventName :(NSDictionary *)params {
  // The bridge eventDispatcher is used to send events from native to JS env
  // No documentation yet on DeviceEventEmitter: https://github.com/facebook/react-native/issues/2819
  [self sendEventWithName: eventName body: params];
}

@end
